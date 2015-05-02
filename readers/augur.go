package readers

import (
	"errors"
	chttp "net/http"
	"net/url"

	"github.com/tyba/oss/sources/social/http"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrPartialResponse      = errors.New("received partial data")
)

var augurInsightsURL = "https://api.augur.io/v2/user"
var augurKey = "2bwn2e88g9dbva8pjolgxeid0nz9m4ne"

type AugurReader struct {
	client *http.Client
}

func NewAugurReader(client *http.Client) *AugurReader {
	return &AugurReader{client}
}

func (a *AugurReader) SearchByEmail(email string) (*AugurInsights, *chttp.Response, error) {
	q := &url.Values{}
	q.Add("email", email)

	r := &AugurInsights{}

	res, err := a.doRequest(q, r)
	if err != nil {
		return nil, res, err
	}

	return r, res, nil
}

func (a *AugurReader) buildURL(q *url.Values) *url.URL {
	q.Add("key", augurKey)

	url, _ := url.Parse(augurInsightsURL)
	url.RawQuery = q.Encode()

	return url
}

func (a *AugurReader) doRequest(q *url.Values, result interface{}) (*chttp.Response, error) {
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
	case 202:
		return res, ErrPartialResponse
	default:
		return res, ErrUnexpectedStatusCode
	}
}

type AugurInsights struct {
	Private        interface{} `json:"PRIVATE"`
	Demographics   interface{} `json:"DEMOGRAPHICS"`
	Psychographics interface{} `json:"PSYCHOGRAPHICS"`
	Geographics    interface{} `json:"GEOGRAPHICS"`
	Profiles       interface{} `json:"PROFILES"`
	Misc           interface{} `json:"MISC"`
	Status         int
}
