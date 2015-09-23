package linkedin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tyba/srcd-domain/models/company"

	"github.com/PuerkitoBio/goquery"
	"github.com/tyba/srcd-rovers/client"
)

const (
	BaseURL      = "https://www.linkedin.com"
	EmployeesURL = BaseURL + "/vsearch/p?f_CC=%d"
)

type LinkedIn struct {
	Cookie string

	client *client.Client
}

func NewLinkedIn(client *client.Client) *LinkedIn {
	return &LinkedIn{client: client}
}

func (g *LinkedIn) GetEmployes(companyId int) (interface{}, error) {
	url := fmt.Sprintf(EmployeesURL, companyId)

	var err error
	var people []person
	for {
		var more []person
		fmt.Printf("Processing %q ...\n", url)
		url, more, err = g.doGetEmployes(url)
		people = append(people, more...)

		if err != nil {
			break
		}

		if url == "" {
			break
		}

	}

	fmt.Printf("Found %d employees\n", len(people))
	fmt.Println(people)

	return nil, err
}

func (g *LinkedIn) doGetEmployes(url string) (
	next string, people []person, err error,
) {
	req, err := client.NewRequest(url)
	if err != nil {
		return
	}

	req.Header.Add("Cookie", g.Cookie)

	doc, res, err := g.client.DoHTML(req)
	if err != nil {
		return
	}

	if res.StatusCode == 404 {
		err = client.NotFound
		return
	}

	return g.parseContent(doc)
}

func (g *LinkedIn) parseContent(doc *goquery.Document) (
	next string, people []person, err error,
) {
	content, err := doc.Find("#voltron_srp_main-content").Html()
	if err != nil {
		return
	}

	//fixing crappy JSON from linkedIn
	//http://stackoverflow.com/questions/30270668/json-loads-giving-exception-that-it-expects-a-value-looks-like-value-is-there
	content = strings.Replace(content, "\\u002d", "-", -1)

	l := len(content)
	js := content[4 : l-3]

	var v voltron
	err = json.Unmarshal([]byte(js), &v)
	if err != nil {
		return
	}

	next = v.getNextURL()
	people = v.getPersons()

	return
}

type voltron struct {
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
						Person person
					}
				}
			} `json:"voltron_unified_search_json"`
		}
	}
}

func (v *voltron) getNextURL() string {
	next := false
	for _, page := range v.Content.Page.V.Search.Data.Pagination.Pages {
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

func (v *voltron) getPersons() []person {
	var o []person
	for _, w := range v.Content.Page.V.Search.Results {
		o = append(o, w.Person)
	}

	return o
}

type person struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Position   string `json:"fmt_headline"`
	LinkedInId int    `json:"id"`
	Location   string `json:"fmt_location"`
}

func (p *person) castToDomain() *company.Employee {
	return &company.Employee{
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		Position:   p.Position,
		LinkedInId: p.LinkedInId,
		Location:   p.Location,
	}
}
