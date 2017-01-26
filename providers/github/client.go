package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	httpTimeout  = 30 * time.Second
	githubApiURL = "https://api.github.com/repositories?since=%d"

	rateLimitLimitHeader     = "X-RateLimit-Limit"
	rateLimitRemainingHeader = "X-RateLimit-Remaining"
	// Link contains the next and first urls of the API endpoint. Example:
	// <https://api.github.com/repositories?since=367>; rel="next", <https://api.github.com/repositories{?since}>; rel="first"
	linkHeader = "Link"
)

type response struct {
	Next         int
	Repositories []*Repository

	Total     int
	Remaining int
}

type client struct {
	c *http.Client
}

func newClient(token string) *client {
	c := &http.Client{}

	if token != "" {
		t := &oauth2.Token{AccessToken: token}
		c = oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(t))
	}

	c.Timeout = httpTimeout

	return &client{c}
}

// Repositories returns a response with the next page id and a list of Repositories.
// It automatically slow down if we are doing requests too fast.
func (c *client) Repositories(since int) (*response, error) {
	start := time.Now()

	u := fmt.Sprintf(githubApiURL, since)
	res, err := c.c.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("request error. Status code %s, %s", res.Status, res.Status)
	}

	repositories, err := c.decode(res.Body)

	total := c.toInt(res.Header.Get(rateLimitLimitHeader))
	remaining := c.toInt(res.Header.Get(rateLimitRemainingHeader))

	minRequestDuration := time.Hour / time.Duration(total)
	defer func() {
		needsWait := minRequestDuration - time.Since(start)
		if needsWait > 0 {
			time.Sleep(needsWait)
		}
	}()

	next := c.next(res)

	return &response{
		Next:         next,
		Repositories: repositories,
		Total:        total,
		Remaining:    remaining,
	}, nil
}

// next parses the HTTP Link response headers and populates the
// various pagination link values in the Response.
func (c *client) next(res *http.Response) int {
	var next int
	if links, ok := res.Header[linkHeader]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")

			// link must at least have href and rel
			if len(segments) < 2 {
				continue
			}

			// ensure href is properly formatted
			if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
				continue
			}

			// try to pull out page parameter
			u, err := url.Parse(segments[0][1 : len(segments[0])-1])
			if err != nil {
				continue
			}
			page := u.Query().Get("page")
			if page == "" {
				continue
			}

			for _, segment := range segments[1:] {
				switch strings.TrimSpace(segment) {
				case `rel="next"`:
					next = c.toInt(page)
				}

			}
		}
	}

	return next
}

func (c *client) toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (c *client) decode(body io.Reader) ([]*Repository, error) {
	var record []*Repository
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}
