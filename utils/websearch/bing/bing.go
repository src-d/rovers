package bing

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)

const (
	timeout   = 30 * time.Second
	keyHeader = "Ocp-Apim-Subscription-Key"

	apiHost   = "api.cognitive.microsoft.com"
	apiPath   = "/bing/v5.0/search"
	apiScheme = "https"

	countParam          = "count"
	offsetParam         = "offset"
	responseFilterParam = "responseFilter"
	queryParam          = "q"

	responseFilterValue = "Webpages"
	countValue          = 50
)

var errQuotaExceeded error = errors.New("Quota exceeded")
var errInvalidKey error = errors.New("Invalid key")
var errTooManyRequests error = errors.New("Too many requests")
var errUnexpected error = errors.New("Bing unexpected error")

type Bing struct {
	apiKey string
	client *http.Client
}

func New(key string) *Bing {
	return &Bing{
		apiKey: key,
		client: &http.Client{Timeout: timeout},
	}
}

func (b *Bing) apiUrl(query string, offset int) *url.URL {
	u := &url.URL{
		Host:   apiHost,
		Path:   apiPath,
		Scheme: apiScheme,
	}

	q := u.Query()
	q.Add(countParam, strconv.Itoa(countValue))
	q.Add(responseFilterParam, responseFilterValue)
	q.Add(queryParam, query)
	q.Add(offsetParam, strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	return u
}

func (b *Bing) newRequest(u *url.URL) *http.Request {
	return &http.Request{
		Header: http.Header{keyHeader: []string{b.apiKey}},
		Method: http.MethodGet,
		URL:    u,
	}
}

func (b *Bing) Search(query string) ([]*url.URL, error) {
	offset := 0
	urls := []*url.URL{}
For:
	for {
		baseUrl := b.apiUrl(query, offset)
		log15.Info("obtaining page for Bing search", "query", query, "offset", offset)
		log15.Debug("actual API URL", "api URL", baseUrl.String())
		resp, err := b.client.Do(b.newRequest(baseUrl))
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode {
		case http.StatusOK:
			result, err := b.getResponse(resp.Body)
			if err != nil {
				return nil, err
			}

			for _, v := range result.WebPages.Values {
				originalURL := v.URL
				resolvedURL, err := b.resolveURL(originalURL)
				if err != nil {
					log15.Error("Error resolving URL, ignoring.", "original URL", originalURL, "error", err)
					continue
				}
				log15.Debug("new resolved URL", "resolved URL", resolvedURL.String())
				urls = append(urls, resolvedURL)
			}
			if b.isLastPage(offset, result.WebPages.TotalEstimatedMatches) {
				break For
			} else {
				offset = offset + countValue
			}
		case http.StatusUnauthorized:
			return nil, errInvalidKey
		case http.StatusForbidden:
			return nil, errQuotaExceeded
		case http.StatusTooManyRequests:
			return nil, errTooManyRequests
		default:
			return nil, errUnexpected
		}
	}

	return urls, nil
}

func (b *Bing) isLastPage(offset int, totalEstimatedMatches int) bool {
	return offset+countValue >= totalEstimatedMatches
}

func (b *Bing) resolveURL(u string) (*url.URL, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	resolvedURL := parsedURL.Query().Get("r")

	return url.Parse(resolvedURL)
}

func (b *Bing) getResponse(body io.ReadCloser) (*BingResult, error) {
	var record BingResult
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}
