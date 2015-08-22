package readers

import (
	"github.com/tyba/srcd-rovers/client"

	"github.com/PuerkitoBio/goquery"
)

var savannahList = "http://savannah.gnu.org/search/?type_of_search=soft&words=%2A&offset=0&max_rows=25000"

type SavannahReader struct {
	client *client.Client
}

func NewSavannahReader(client *client.Client) *SavannahReader {
	return &SavannahReader{client}
}

func (s *SavannahReader) GetRepositories() (*SavannahResult, error) {
	req, err := client.NewRequest(savannahList)
	if err != nil {
		return nil, err
	}

	doc, _, err := s.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	return NewSavannahResult(doc), nil
}

type SavannahRepository struct {
	Project     string
	Description string
	Type        string
	Location    string
	Link        string
}

type SavannahResult struct {
	Results []SavannahRepository
}

func NewSavannahResult(doc *goquery.Document) *SavannahResult {
	res := &SavannahResult{Results: make([]SavannahRepository, 0)}

	doc.Find(".box").Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() != 3 {
			return
		}

		r := SavannahRepository{}
		r.Project = tds.Slice(0, 1).Text()
		r.Description = tds.Slice(1, 2).Text()
		r.Type = tds.Slice(2, 3).Text()
		r.Link, _ = s.Find("a").Attr("href")

		res.Results = append(res.Results, r)
	})

	return res
}
