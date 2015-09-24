package linkedin

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/tyba/srcd-domain/models/company"
	"github.com/tyba/srcd-rovers/client"

	"github.com/PuerkitoBio/goquery"
)

const (
	BaseURL      = "https://www.linkedin.com"
	EmployeesURL = BaseURL + "/vsearch/p?f_CC=%d"
)

type LinkedInWebCrawler struct {
	client *client.Client
	cookie string
}

func NewLinkedInWebCrawler(client *client.Client, cookie string) *LinkedInWebCrawler {
	return &LinkedInWebCrawler{client: client, cookie: cookie}
}

func (li *LinkedInWebCrawler) GetEmployees(companyId int) (
	people []company.Employee, err error,
) {
	start := time.Now()
	url := fmt.Sprintf(EmployeesURL, companyId)

	for {
		var more []Person
		log15.Info("Processing", "url", url)
		url, more, err = li.doGetEmployes(url)

		for _, person := range more {
			people = append(people, person.ToDomainCompanyEmployee())
		}

		if err != nil || url == "" {
			break
		}
	}

	log15.Info("Done",
		"elapsed", time.Since(start),
		"found", len(people),
	)
	for idx, person := range people {
		log15.Debug("Person", "idx", idx, "person", person)
	}
	return people, err
}

func (li *LinkedInWebCrawler) doGetEmployes(url string) (
	next string, people []Person, err error,
) {
	req, err := client.NewRequest(url)
	if err != nil {
		return
	}
	req.Header.Add("Cookie", li.cookie)

	doc, res, err := li.client.DoHTML(req)
	if err != nil {
		return
	}
	if res.StatusCode == 404 {
		err = client.NotFound
		return
	}
	return li.parseContent(doc)
}

func (li *LinkedInWebCrawler) parseContent(doc *goquery.Document) (
	next string, people []Person, err error,
) {
	content, err := doc.Find("#voltron_srp_main-content").Html()
	if err != nil {
		return
	}

	// Fix encoding issues with LinkedIn's JSON:
	// Source: http://stackoverflow.com/q/30270668
	content = strings.Replace(content, "\\u002d", "-", -1)

	length := len(content)
	jsonBlob := content[4 : length-3]

	var data LinkedInData
	err = json.Unmarshal([]byte(jsonBlob), &data)
	if err != nil {
		return
	}

	next = data.getNextURL()
	people = data.getPeople()
	return
}

// fat ass LinkedIn format
type LinkedInData struct {
	Content struct {
		Page struct {
			V struct {
				Search struct {
					Data struct {
						Pagination struct {
							Pages []struct {
								Current bool   `json:"isCurrentPage"`
								URL     string `json:"pageURL"`
							}
						} `json:"resultPagination"`
					} `json:"baseData"`
					Results []struct {
						Person Person
					}
				}
			} `json:"voltron_unified_search_json"`
		}
	}
}

func (lid *LinkedInData) getNextURL() string {
	next := false
	for _, page := range lid.Content.Page.V.Search.Data.Pagination.Pages {
		if page.Current {
			next = true
			continue
		}

		if next {
			return BaseURL + page.URL
		}
	}

	return ""
}

func (lid *LinkedInData) getPeople() []Person {
	var people []Person
	for _, result := range lid.Content.Page.V.Search.Results {
		people = append(people, result.Person)
	}
	return people
}

type Person struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	LinkedInId int    `json:"id"`
	Location   string `json:"fmt_location"`
	Position   string `json:"fmt_headline"`
}

func (p *Person) ToDomainCompanyEmployee() company.Employee {
	return company.Employee{
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		LinkedInId: p.LinkedInId,
		Location:   p.Location,
		Position:   p.Position,
	}
}
