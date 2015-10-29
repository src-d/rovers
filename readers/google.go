package readers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/src-d/rovers/client"

	"github.com/PuerkitoBio/goquery"
)

const googleSearch = "http://www.google.com/search?hl=en&q=%s&ie=UTF-8&btnG=Google+Search&inurl=https"

type GoogleSearchReader struct {
	client *client.Client
}

func NewGoogleSearchReader(client *client.Client) *GoogleSearchReader {
	return &GoogleSearchReader{client}
}

func (g *GoogleSearchReader) SearchByQuery(query string) (*GoogleSearchResult, error) {
	query = strings.Replace(query, " ", "+", -1)
	req, err := client.NewRequest(fmt.Sprintf(googleSearch, query))
	if err != nil {
		return nil, err
	}

	doc, _, err := g.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	return NewGoogleSearchResult(doc), nil
}

type GoogleSearchResult struct {
	Search  string
	Results []result
}

func NewGoogleSearchResult(doc *goquery.Document) *GoogleSearchResult {
	g := &GoogleSearchResult{}

	g.Search, _ = doc.Find(".lst").Attr("value")

	g.Results = make([]result, 0)
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		link, _ := getURLFromResult(href)

		g.Results = append(g.Results, result{
			Name: s.Text(),
			Link: link,
		})
	})

	return g
}

func (g *GoogleSearchResult) FindByHost(host string) []result {
	results := make([]result, 0)
	for _, result := range g.Results {
		if result.MatchHost(host) {
			results = append(results, result)
		}
	}

	return results
}

type result struct {
	Name string
	Link string
}

func (r *result) GetHost() string {
	if r.Link == "" {
		return ""
	}

	u, _ := url.Parse(r.Link)
	return u.Host
}

func (r *result) MatchHost(host string) bool {
	if r.GetHost() == host {
		return true
	}

	if strings.HasSuffix(r.GetHost(), host) {
		return true
	}

	return false
}

func getURLFromResult(href string) (string, error) {
	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return u.Query().Get("q"), nil
}
