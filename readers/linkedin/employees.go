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
	UserAgent                  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:41.0) Gecko/20100101 Firefox/41.0"
	CookieFixture              = `bcookie="v=2&21733669-d772-42c2-8d0c-6631f3b494b7"; bscookie="v=1&20150904164640a97c0766-597a-4097-83e6-aa8ddf3ce0f6AQEBd5soPhRmh7HBB1afmID0u3OdZN6_"; visit="v=1&M"; sessionid="eyJkamFuZ29fdGltZXpvbmUiOiJFdXJvcGUvQmVybGluIn0:1ZbNU6:6Qlqhdr9bWBDhQNGA567O14PNrY"; csrftoken=3T0lfgzrejTNI80YLJduCviqdTtqYxLx; __utma=23068709.221348018.1441385201.1442832110.1442832110.1; __utmz=23068709.1442832110.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmv=23068709.guest; L1c=38153dc8; wutan=5XW4mxShwpmVW0IhvylSv8h5tXLFDtAg9qc5/+mAQoI=; L1e=1e953b94; li_at=AQEDAQD4rTcDfJngAAABT8gc3_wAAAFP-euXSk4ArFmrYrEtDQ5abVZGsQQ9cLzu1Htm0WZ22vJ5bhJRMpd1-o9a6FET44xN5vG90I5Mst_NtpsPnGUyNPcJR-O3sT15_SXimb7ObEwtDpqiISBkAjCz; liap=true; sl="v=1&x_gci"; JSESSIONID="ajax:7765494870178832670"; oz_props_fetch_size1_16297271=6; share_setting=PUBLIC; sdsc=1%3A1SZM1shxDNbLt36wZwCgPgvN58iw%3D; lidc="b=TB71:g=114:u=115:i=1443005316:t=1443074927:s=AQE_9GpXb_1hvkwazhni2mcwgIO3bXA4"; _ga=GA1.2.221348018.1441385201; _gat=1; RT=s=1443005316840&r=https%3A%2F%2Fwww.linkedin.com%2Fcompany%2F924688; _lipt=0_0DWUvqOlLiwrAUp1_qzuTvYKhO30OcJ9TEc3PczGd9dDgjcZ3KAXmhKA8eI6zryVYkmWcw-jFWWKI9Y1axh_16jv7p-SSo-G8o3kNVxDF5uQ_pKPUcwECfQt5cKUp1RtrgPEwPNCNvfZj_EL8mSAkNMgC_n_MQ_djS9R8jd-4DNXjh6uHbexZL3ZMiyEwiWXVvkjSDpXz2ZY9mPD022h5J6eQtIt313hKiRyrB3sMn_T6I0bxXRmS3Ob-Q3TW_JFl3pU9euQDHhmsn3eaOVcHf; lang="v=2&lang=en-us"`
	CookieFixture2             = `lang="v=2&lang=en-us"; bcookie="v=2&234c28f0-6149-4df1-8db8-e00e8c86dbfb"; bscookie="v=1&20151009141721b2611d93-ff0b-46ac-8d31-1ff263045393AQEJmqishE1t6pWXplq4MXqe8RF_igAy"; lidc="b=TB30:g=256:u=260:i=1444400255:t=1444436333:s=AQEepfx6DomMdzYs_qoESCYJ_KhWOhXc"; visit="v=1&M"; sl="v=1&6jhXN"; liap=true; li_at=AQEDAQB8ujIEstvuAAABUEz2z4cAAAFQTWSsh04AE0AKeyq1AV1FriAShK-FieYa8DAY68EopV3Y6jqwDpJbE7MNLyfC6vGB4s0zrLh29q3PCCjAaMS07n3Mk7RNwFNhI0JMDS-19ysfUzB-O6BylgvJ; RT=s=1444400254669&r=https%3A%2F%2Fwww.linkedin.com%2F; _lipt=0_cXjzxyYyIs1sRlATNKJxy5F4ZGSV6vEuSl8B1qdmoXy_W0iUDD6xAsfhLR3vflzmFc2V2m5-pLDhKwIEP7KYQBAErt9i3raE6Led1OIGfnS2Z1_9seJ-KwtJVE1D392LFKdK1D_e7PKNnrmE_5jIoFCm-zJf2jIzHCEYSMwsXbtA5O_EnpKnHxuHh5fzSXmQVHOXtIL-icI_jeGoEfEwSMTU6ui08Fo1IQxzrCkSW7n_NJEEf3rHf1e7eOV7HFHoXZOibZaV0owY4xBMMwsrR_AOgoHjlpeyFhjXyHM4FzeqLrMsUz-ozGMyuk0JUCQt6xv7ZQ5Eqri6dh4aWRvSMepYnY3QzqsrLcHJphtSjia6eQtIt313hKiRyrB3sMn_S_OCBSbwDtc1Vwi3NncuhvlJCd7weCShFpjj1NEJGs8; L1c=5824ada7; oz_props_fetch_size1_8174130=15; wutan=AtY6lFc9IWIpt31Z2cwh3ndwUFzpOw8Nr4vLzBG1g3M=; share_setting=PUBLIC; sdsc=22%3A1%2C1444400085057%7ECAOR%2C0uPkgz7kt9ua2ENuBx8vaxzEmF00%3D`
	BaseURL                    = "https://www.linkedin.com"
	EmployeesURL               = BaseURL + "/vsearch/p?f_CC=%d"
	LinkedInEmployeesRateLimit = 5 * time.Second
)

type LinkedInWebCrawler struct {
	client *client.Client
	cookie string
}

func NewLinkedInWebCrawler(client *client.Client, cookie string) *LinkedInWebCrawler {
	switch cookie {
	case "fixture":
		cookie = CookieFixture
	case "fixture2":
		cookie = CookieFixture2
	case "":
		panic("empty cookie")
	}
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
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Cookie", li.cookie)

	doc, res, err := li.client.DoHTML(req)
	if err != nil {
		return
	}
	log15.Debug("DoHTML", "url", req.URL, "status", res.StatusCode)
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
