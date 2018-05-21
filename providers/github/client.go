package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/inconshreveable/log15.v2"

	"github.com/src-d/rovers/providers/github/model"
)

const (
	httpTimeout  = 30 * time.Second
	githubApiURL = "https://api.github.com/repositories?since=%d"

	rateLimitLimitHeader     = "X-RateLimit-Limit"
	rateLimitRemainingHeader = "X-RateLimit-Remaining"
	rateLimitResetHeader     = "X-RateLimit-Reset"
)

type response struct {
	Next         int
	Repositories []*model.Repository
}

type errorResponse struct {
	Message string `json:"message"`
}

type client struct {
	c        *http.Client
	endpoint string
}

func newClient(token string) *client {
	c := &http.Client{}

	if token != "" {
		t := &oauth2.Token{AccessToken: token}
		c = oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(t))
	}

	c.Timeout = httpTimeout

	return &client{c, githubApiURL}
}

// Repositories returns a response with the next page id and a list of Repositories.
// It automatically slow down if we are doing requests too fast.
func (c *client) Repositories(since int) (*response, error) {
	for {
		resp, retry, err := c.repositories(since)
		if !retry {
			return resp, err
		}

		log15.Warn("got retryable error", "err", err)
	}
}

func (c *client) repositories(since int) (*response, bool, error) {
	u := fmt.Sprintf(c.endpoint, since)
	res, err := c.c.Get(u)
	if err != nil {
		return nil, false, err
	}

	defer c.wait(res)
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, c.isRateLimitError(res), c.decodeError(res)
	}

	repositories, err := c.decode(res.Body)
	if err != nil {
		return nil, false, err
	}

	// remove those repositories that GitHub API encoded as null
	// in the JSON reponse and were decoded as a nil element in the
	// *model.Repository slice.
	validRepositories := make([]*model.Repository, 0, len(repositories))
	for _, repo := range repositories {
		if repo != nil {
			validRepositories = append(validRepositories, repo)
		}
	}

	next := 0
	if len(validRepositories) != 0 {
		next = validRepositories[len(validRepositories)-1].GithubID
	}

	return &response{
		Next:         next,
		Repositories: validRepositories,
	}, false, nil
}

func (c *client) toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (c *client) wait(res *http.Response) {
	if res.Header.Get(rateLimitRemainingHeader) == "" {
		return
	}

	remaining := c.toInt(res.Header.Get(rateLimitRemainingHeader))

	now := time.Now().UTC().Unix()
	resetTime := int64(c.toInt(res.Header.Get(rateLimitResetHeader)))
	timeToReset := time.Duration(resetTime-now) * time.Second
	if timeToReset < 0 || timeToReset > 1*time.Hour {
		// If this happens, the system clock is probably wrong, so we assume we
		// are at the beginning of the window and consider only total requests
		// per hour.
		timeToReset = 1 * time.Hour
		remaining = c.toInt(res.Header.Get(rateLimitLimitHeader))
	}

	waitTime := timeToReset / time.Duration(remaining+1)
	if waitTime > 0 {
		log15.Debug("waiting", "wait", waitTime)
		time.Sleep(waitTime)
	}
}

func (c *client) isRateLimitError(res *http.Response) bool {
	return res.StatusCode == 403 &&
		res.Header.Get(rateLimitRemainingHeader) == "0"
}

func (c *client) decode(body io.Reader) ([]*model.Repository, error) {
	var record []*model.Repository
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}

func (c *client) decodeError(resp *http.Response) error {
	var errResp = &errorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(errResp); err != nil {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return fmt.Errorf("HTTP error: %s, %s", resp.Status, errResp.Message)
}
