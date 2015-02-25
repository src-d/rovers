package sources

import (
	chttp "net/http"
	"net/url"

	"github.com/tyba/opensource-search/sources/social/http"
)

var bitbucketURL = "https://api.bitbucket.org/2.0/repositories"

type Bitbucket struct {
	client *http.Client
}

func NewBitbucket(client *http.Client) *Bitbucket {
	return &Bitbucket{client}
}

func (a *Bitbucket) GetRepositories(q url.Values) (*BitbucketPagedResult, error) {
	r := &BitbucketPagedResult{}

	_, err := a.doRequest(q, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (a *Bitbucket) buildURL(q url.Values) *url.URL {
	url, _ := url.Parse(bitbucketURL)
	if q.Get("page") != "" {
		url.RawQuery = q.Encode()
	}

	return url
}

func (a *Bitbucket) doRequest(q url.Values, result interface{}) (*chttp.Response, error) {
	req, err := http.NewRequest(a.buildURL(q).String())
	if err != nil {
		return nil, err
	}

	res, err := a.client.DoJSON(req, result)
	if err != nil {
		return res, err
	}

	switch res.StatusCode {
	case 200:
		return res, nil
	default:
		return res, ErrUnexpectedStatusCode
	}
}

type BitbucketPagedResult struct {
	Page       int           `json:"page"`
	PageLength int           `json:"pagelen"`
	Values     []interface{} `json:"values"`
	Next       *URL          `json:"next"`
}

type URL struct {
	*url.URL
}

func (u *URL) MarshalJSON() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *URL) UnmarshalJSON(b []byte) error {
	o, err := url.Parse(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	u.URL = o

	return nil
}
