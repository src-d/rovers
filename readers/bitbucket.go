package readers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/tyba/srcd-rovers/client"
)

var bitbucketURL = "https://api.bitbucket.org/2.0/repositories"

type BitbucketAPI struct {
	client *client.Client
}

func NewBitbucketAPI(client *client.Client) *BitbucketAPI {
	return &BitbucketAPI{client}
}

func (a *BitbucketAPI) GetRepositories(q url.Values) (*BitbucketPagedResult, error) {
	r := &BitbucketPagedResult{}

	_, err := a.doRequest(q, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (a *BitbucketAPI) buildURL(q url.Values) *url.URL {
	u, _ := url.Parse(bitbucketURL)
	if q.Get("page") != "" {
		u.RawQuery = q.Encode()
	}

	return u
}

func (a *BitbucketAPI) doRequest(q url.Values, result interface{}) (*http.Response, error) {
	req, err := client.NewRequest(a.buildURL(q).String())
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
	Page       int          `json:"page"`
	PageLength int          `json:"pagelen"`
	Values     []Repository `json:"values"`
	Next       *URL         `json:"next"`
}

type Repository struct {
	CreatedOn   time.Time `json:"created_on"`
	Description string    `json:"description"`
	ForkPolicy  string    `json:"fork_policy"`
	FullName    string    `json:"full_name"`
	HasIssues   bool      `json:"has_issues"`
	HasWiki     bool      `json:"has_wiki"`
	IsPrivate   bool      `json:"is_private"`
	Language    string    `json:"language"`
	Links       struct {
		Avatar Href `json:"avatar"`
		Clone  []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
		Html Href `json:"html"`
		Self Href `json:"self"`
	} `json:"links"`
	Name  string `json:"name"`
	Owner struct {
		DisplayName string `json:"display_name"`
		Links       struct {
			Avatar Href `json:"avatar"`
			Html   Href `json:"html"`
			Self   Href `json:"self"`
		} `json:"links"`
		Type     string `json:"type"`
		Username string `json:"username"`
		UUID     string `json:"uuid"`
	} `json:"owner"`
	SCM       string    `json:"scm"`
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	UpdatedOn time.Time `json:"updated_on"`
	UUID      string    `json:"uuid"`
}

type Href struct {
	Href string `json:"href"`
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
