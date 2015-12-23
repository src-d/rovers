package linkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/src-d/rovers/client"
	"gop.kg/src-d/domain@v2.1/models/company"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	UserAgent                  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:41.0) Gecko/20100101 Firefox/41.0"
	CookieFixtureEiso          = `JSESSIONID="ajax:7889397718783924060"; lang="v=2&lang=en-us"; bcookie="v=2&3bbf145d-e36b-4268-84c0-f8cf4977dadd"; bscookie="v=1&20151223232515075d4c39-8e12-459c-8b76-02cb415ae84bAQHkBnswRJeLGS1Af__oX2QPMLMPbNFK"; lidc="b=TB30:g=320:u=307:i=1450913129:t=1450991126:s=AQHhT5PQpc9-i8g0iMfiEsXrbuozMM6c"; _ga=GA1.2.811094040.1450913116; _gat=1; sl="v=1&iW-tB"; visit="v=1&M"; li_at=AQEDAQB8ujICkmAoAAABUdEpZJoAAAFR0ZdBmk4AWOcioswTS1skMDQ_0CpSsulJnK_o7s8_7d1K0ZzRzb-Cf0qzgBaExdzWKIPDm1Fc4GZD2B2eYsZatOKclz0-JqNOWTQPupolqw8hP4nb3I5LNuWu; liap=true; RT=s=1450913130316&r=https%3A%2F%2Fwww.linkedin.com%2F; _lipt=0_6Y2YMwAcAPLi2pqIJ3KivaKqW8R0GSzU-1ibPUlEe99b12fq1cTiLEGaAjBOsTL2hG_oTs75skJpZ3Xtb-Hvas1zA5rfxvsaci-f3DO8c3wJXVVMwnnFrfb0KZC7M-TLRhWSRIYLbhPFwKC65Wfk28w6i_yU7kIeb5P9bwf3-7KY27SWB5mtUllB2aDfI7XDRoPdw-QVGzvz5DgsRDI1gLTG6Q6pIDis0z2WvqvSlUHcQVHJvwlv8ozEBGI4Apo-dj2zxXuq10ZuGZgZXIW-MBhL-vaxNHO_3fUHDmkLoLypluRyrIa4qjjL2p0hZpn8Sxv7ZQ5Eqri6dh4aWRvSMepYnY3QzqsrLcHJphtSjia6eQtIt313hKiRyrB3sMn_S_OCBSbwDtc1Vwi3NncuhvlJCd7weCShFpjj1NEJGs8; L1c=75c8dfea; oz_props_fetch_size1_8174130=15; wutan=ELXkTYwXemKkvWz9g7hQbXP/TUjU83axuSH1p+aoAN8=; share_setting=PUBLIC; sdsc=1%3A1SZM1shxDNbLt36wZwCgPgvN58iw%3D`
	BaseURL                    = "https://www.linkedin.com"
	EmployeesURL               = BaseURL + "/vsearch/p?f_CC=%d"
	LinkedInEmployeesRateLimit = 5 * time.Second
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

		if err != nil {
			log15.Error("GetEmployees", "error", err)
			break
		}
		if url == "" {
			break
		}
	}

	log15.Info("Done",
		"elapsed", time.Since(start),
		"found", len(people),
	)
	// for idx, person := range people {
	// 	log15.Debug("Person", "idx", idx, "person", person)
	// }
	return people, err
}

func (li *LinkedInWebCrawler) doGetEmployes(url string) (
	next string, people []Person, err error,
) {
	start := time.Now()
	defer func() {
		needsWait := LinkedInEmployeesRateLimit - time.Since(start)
		if needsWait > 0 {
			log15.Debug("Waiting", "duration", needsWait)
			time.Sleep(needsWait)
		}
	}()

	req, err := client.NewRequest(url)
	if err != nil {
		return
	}

	req.Header.Set("Cookie", li.cookie)

	res, err := li.client.Do(req)
	if err != nil {
		return
	}
	log15.Debug("Do", "url", req.URL, "status", res.StatusCode)
	if res.StatusCode == 404 {
		err = client.NotFound
		return
	}

	doc, err := li.preprocessContent(res)
	if err != nil {
		return
	}

	return li.parseContent(doc)
}

// goquery will transform `&quot;` into `"` even if it's inside a HTML comment
// so we need to replace all of those first by some non-harming character first,
// like `'`, so we can JSON decode the `Voltron` payload succesfully
func (l *LinkedInWebCrawler) preprocessContent(res *http.Response) (*goquery.Document, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	idx := bytes.Index(body, []byte("voltron_srp_main-content"))
	if idx > -1 {
		log15.Info("FOUND voltron payload")
	} else {
		log15.Info("NOT FOUND voltron payload")
	}

	body = bytes.Replace(body, []byte("&quot;"), []byte(`\"`), -1)

	reader := bytes.NewBuffer(body)
	return goquery.NewDocumentFromReader(reader)
}

func (li *LinkedInWebCrawler) parseContent(doc *goquery.Document) (
	next string, people []Person, err error,
) {
	content, err := doc.Find("#voltron_srp_main-content").Html()
	if err != nil {
		return
	}

	// Fix encoding issues with LinkedIn's JSON:
	content = strings.TrimPrefix(content, "<!--")
	content = strings.TrimSuffix(content, "-->")
	content = strings.Replace(content, `\u003c`, "<", -1)
	content = strings.Replace(content, `\u003e`, ">", -1)
	content = strings.Replace(content, `\u002d`, "-", -1)

	if len(content) == 0 {
		log15.Warn("No JSON found for page")
		return
	}

	var data LinkedInData
	contentBytes := []byte(content)
	err = json.Unmarshal(contentBytes, &data)
	if err != nil {
		log15.Error("json.Unmarshal", "error", err)
		if serr, ok := err.(*json.SyntaxError); ok {
			log.Println("SyntaxError at offset:", serr.Offset)
			log.Printf("%s\n", contentBytes[serr.Offset-20:serr.Offset+20])
		}
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
