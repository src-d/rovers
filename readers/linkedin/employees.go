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
	"gop.kg/src-d/domain@v6/models/company"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	UserAgent                  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:41.0) Gecko/20100101 Firefox/41.0"
	CookieFixtureEiso          = `bcookie="v=2&a6df69ca-69d9-43e5-856b-0b77adc95a04"; bscookie="v=1&2016012016431869318f48-9ac1-4d89-82cd-997e7f977e54AQFdmcU0uOQ2UUtCLEOL20ZDxELWvPXs"; _ga=GA1.2.544046633.1453308203; visit="v=1&M"; JSESSIONID="ajax:6685549765576733025"; liap=true; li_at=AQEDAQB8ujIFVd69AAABVJnnixYAAAFUmlVoFksA0jD8N8qow8FeFCIZOej2iFseI9Yebx8e6xY2MSavthC6Da9FrUB_jQ_2dLZ8Ds4rvfD9p6I9u7ZoWNwwRETsDhKjmQoeeRA7iUW72UbgmXBU4YEC; _lipt=0_1Mq7gN3oKU8DI3XUemApNB3grY_sQCmrK8obxyc2KyN1JwAiSFHmeStoZYWYu_iPQzJD1nk4lI4tyUSSVtK_c1chlz2TUPHUClox7wHfNlrhq1dYOgod9L7rukf2WoG9XmJctL_ah4e-GEHWyHx1sdzNot-1_HIhjje9Mfn5o8rC8ZcriXrRTQxYHKJB71f18NNZcudFBOhpit8keKqGx2UcKclJILQoNBGNB-TAYp8Z3JCPJZSdniKyMkAtlYftSxv7ZQ5Eqri6dh4aWRvSMc7822GdbKqY2O7irGJuVbm2oQuL1fosyq99AjG7nhX2ld12C7t1jS-6cHaYVdD5d9nidupkSCHrutV6zt2gw7H; lang="v=2&lang=en-us"; lidc="b=LB30:g=518:u=387:i=1462870963:t=1462946397:s=AQH7JaeF7NmVFxf7xo98Nwd4WqScbOUs"; _gat=1; L1e=a98019b; sl="v=1&H96VN"; RT=s=1462870969586&r=https%3A%2F%2Fwww.linkedin.com%2F; L1c=5eba1285; oz_props_fetch_size1_8174130=15; wutan=k12JhcE/r/rpOzWbnskhsZGyc+eEsONK/zRksDDqRKo=; share_setting=PUBLIC; sdsc=1%3A1SZM1shxDNbLt36wZwCgPgvN58iw%3D`
	BaseURL                    = "https://www.linkedin.com"
	EmployeesURL               = BaseURL + "/vsearch/p"
	IdFilter                   = "?f_CC=%d"
	TitleFilter                = "&title=%s&titleScope=CP"
	KeywordJoiner              = "%20OR%20"
	LinkedInEmployeesRateLimit = 5 * time.Second
)

var (
	Keywords = []string{
		"architect",
		"chief",
		"coder",
		"cto",
		"dataops",
		"desarrollador",
		"developer",
		"devops",
		"engineer",
		"engineering",
		"programador",
		"programmer",
		"software",
		"system",
		"systems",
	}
	JoinedTitles = strings.Join(Keywords, KeywordJoiner)
	Titles       = fmt.Sprintf(TitleFilter, JoinedTitles)
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

	url := EmployeesURL + fmt.Sprintf(IdFilter, companyId) + Titles

	for {
		var more []Person
		// log15.Info("Processing", "url", url)
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
		log15.Info("FOUND voltron payload", "url", res.Request.URL)
	} else {
		log15.Info("NOT FOUND voltron payload", "url", res.Request.URL)
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
