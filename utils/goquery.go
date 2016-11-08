package utils

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const httpTimeout = 30 * time.Second

// Goquery client used to scrape data from web pages.
// This struct is necessary because the default timeout
// into http.DefaultClient is 0. Because of that, we use
// a custom instance of http.Client.
type GoqueryClient struct {
	client http.Client
}

func NewDefaultGoqueryClient() *GoqueryClient {
	return &GoqueryClient{
		client: http.Client{
			Timeout: httpTimeout,
		},
	}
}

// Generates a new Goquery document from the given url using a custom http.Client, if possible.
func (gq *GoqueryClient) NewDocument(url string) (*goquery.Document, error) {
	resp, err := gq.client.Get(url)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(resp)
}
