package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/src-d/rovers/providers/bitbucket/models"
)

const (
	httpTimeout = 30 * time.Second

	baseURL          = "https://api.bitbucket.org/2.0/"
	repositoriesPath = "repositories"

	afterParam   = "after"
	pagelenParam = "pagelen"

	pagelenValue = 100
)

type response struct {
	Pagelen      int                 `json:"pagelen"`
	Repositories models.Repositories `json:"values"`
	Next         string              `json:"next"`
}

type client struct {
	c *http.Client
}

func newClient() *client {
	return &client{
		c: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (c *client) parse(after string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, repositoriesPath)

	q := u.Query()
	if after != "" {
		q.Add(afterParam, after)
	}
	q.Add(pagelenParam, strconv.Itoa(pagelenValue))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *client) decode(body io.Reader) (*response, error) {
	var record response
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}

	if record.Next != "" {
		u, err := url.Parse(record.Next)
		if err != nil {
			return nil, err
		}
		record.Next = u.Query().Get(afterParam)
	}

	return &record, nil
}

// Repositories returns Bitbucket API response with 100 repositories max.
// 'after' is a string in a very specific date time format used to get data from a specific time.
// If you want to get the first page, use "". Every response result has a 'Next' field.
// It could be used to get the next page calling again to Repositories method.
// If you try to obtain a page that doesn't exists, error 'io.EOF' is returned.
func (c *client) Repositories(after string) (*response, error) {
	u, err := c.parse(after)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.c.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	response, err := c.decode(res.Body)
	if err != nil {
		return nil, err
	}

	if len(response.Repositories) == 0 {
		return nil, io.EOF
	}

	return response, nil
}
