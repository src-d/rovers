package sources

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/tyba/opensource-search/sources/social/http"

	"github.com/PuerkitoBio/goquery"
)

type LinkedIn struct {
	client *http.Client
}

func NewLinkedIn(client *http.Client) *LinkedIn {
	return &LinkedIn{client}
}

func (l *LinkedIn) GetProfileByURL(url string) (*LinkedInProfile, error) {
	req, err := http.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, res, err := l.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 error code: %d", res.StatusCode)
	}

	return NewLinkedInProfile(url, doc), nil
}

type LinkedInProfile struct {
	Created      time.Time
	URL          string
	Connections  int
	FullName     string
	Title        string
	Locality     string
	Industry     string
	Websites     map[string]string
	Current      map[string]string
	Previous     map[string]string
	Summary      string
	Experience   []experience
	Languages    map[string]string
	Skills       []string
	Education    []education
	Patents      []patent
	Projects     []project
	Publications []publication
	Related      []person
}

type experience struct {
	Title     string
	Company   string
	Link      string
	StartDate time.Time
	EndDate   time.Time
	Locality  string
	Summary   string
}

type education struct {
	Summary string
	Degree  string
	Major   string
	Extra   string
	Date    string
}

type patent struct {
	Title   string
	Date    string
	Id      string
	Summary string
	Persons []person
}

type project struct {
	Title   string
	Link    string
	Date    string
	Summary string
	Persons []person
}

type publication struct {
	Title   string
	Link    string
	Date    string
	Summary string
	Persons []person
}

type person struct {
	Name string
	Link string
}

func NewLinkedInProfile(url string, doc *goquery.Document) *LinkedInProfile {
	p := &LinkedInProfile{URL: url, Created: time.Now()}
	fillConnections(p, doc)
	fillBasicInfo(p, doc)
	fillExperience(p, doc)
	fillLanguages(p, doc)
	fillSkills(p, doc)
	fillEducation(p, doc)
	fillPatents(p, doc)
	fillProjects(p, doc)
	fillPublications(p, doc)
	fillAlsoViewed(p, doc)

	return p
}

func fillConnections(p *LinkedInProfile, doc *goquery.Document) {
	c := doc.Find(".profile-overview .member-connections strong").Text()
	if c == "500+" {
		p.Connections = 500
	} else {
		p.Connections, _ = strconv.Atoi(c)
	}
}

func fillBasicInfo(p *LinkedInProfile, doc *goquery.Document) {
	p.FullName = doc.Find(".full-name").Text()
	p.Title = doc.Find(".title").Text()
	p.Industry = doc.Find("#location .industry").Text()
	p.Locality = doc.Find("#location .locality").Text()

	p.Websites = make(map[string]string, 0)
	doc.Find("#overview-summary-websites").Find("li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		url, _ := getURLFromRedirect(href)

		p.Websites[s.Text()] = url
	})

	p.Current = make(map[string]string, 0)
	doc.Find("#overview-summary-current").Find("li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")

		p.Current[s.Text()] = href
	})

	p.Previous = make(map[string]string, 0)
	doc.Find("#overview-summary-past").Find("li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")

		p.Previous[s.Text()] = href
	})

	p.Summary = doc.Find("#summary-item").Text()
}

func fillExperience(p *LinkedInProfile, doc *goquery.Document) {
	p.Experience = make([]experience, 0)
	doc.Find("#background-experience").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		e := experience{}
		e.Title = s.Find("h4").Text()
		e.Company = s.Find("h5").Text()
		e.Link, _ = s.Find("h5").Find("a").Attr("href")
		e.Locality = s.Find(".locality").Text()
		e.Summary = s.Find(".description").Text()

		s.Find("time").Each(func(i int, s *goquery.Selection) {
			date, _ := getTimeFromString(s.Text())
			if i == 0 {
				e.StartDate = date
			} else {
				e.EndDate = date
			}
		})

		p.Experience = append(p.Experience, e)
	})
}

func fillLanguages(p *LinkedInProfile, doc *goquery.Document) {
	p.Languages = make(map[string]string, 0)
	doc.Find("#languages").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		language := s.Find("h4").Text()
		level := s.Find(".languages-proficiency").Text()

		p.Languages[language] = level
	})
}

func fillSkills(p *LinkedInProfile, doc *goquery.Document) {
	p.Skills = make([]string, 0)
	doc.Find(".endorse-item ").Each(func(i int, s *goquery.Selection) {
		if _, ok := s.Attr("id"); !ok {
			p.Skills = append(p.Skills, s.Text())
		}
	})
}

func fillEducation(p *LinkedInProfile, doc *goquery.Document) {
	p.Education = make([]education, 0)
	doc.Find("#background-education").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		e := education{}
		e.Summary = s.Find("h4").Text()
		e.Degree = s.Find(".degree").Text()
		e.Major = s.Find(".major").Text()
		e.Date = s.Find(".education-date").Text()
		e.Extra = s.Find("h5").Text()

		p.Education = append(p.Education, e)
	})
}

func fillPatents(p *LinkedInProfile, doc *goquery.Document) {
	p.Patents = make([]patent, 0)
	doc.Find("#background-patents").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		pa := patent{}
		pa.Title = s.Find("h4 span:first-of-type").Text()
		pa.Id = s.Find("h5").Text()
		pa.Date = s.Find(".patents-date").Text()
		pa.Summary = s.Find(".description").Text()
		pa.Persons = getPersons(s.Find(".associated-endorsements"))

		p.Patents = append(p.Patents, pa)
	})
}

func fillProjects(p *LinkedInProfile, doc *goquery.Document) {
	p.Projects = make([]project, 0)
	doc.Find("#background-projects").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		pr := project{}
		pr.Title = s.Find("h4 span:first-of-type").Text()
		href, _ := s.Find("h4").Find("a").Attr("href")
		pr.Link, _ = getURLFromRedirect(href)
		pr.Date = s.Find("time").Text()
		pr.Summary = s.Find(".description").Text()
		pr.Persons = getPersons(s.Find(".associated-endorsements"))

		p.Projects = append(p.Projects, pr)
	})
}

func fillPublications(p *LinkedInProfile, doc *goquery.Document) {
	p.Publications = make([]publication, 0)
	doc.Find("#background-publications").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		pu := publication{}
		pu.Title = s.Find("h4").Text()
		href, _ := s.Find("h4").Find("a").Attr("href")
		pu.Link, _ = getURLFromRedirect(href)
		pu.Date = s.Find(".publication-date").Text()
		pu.Summary = s.Find(".description").Text()
		pu.Persons = getPersons(s.Find(".associated-endorsements"))

		p.Publications = append(p.Publications, pu)
	})
}

func fillAlsoViewed(p *LinkedInProfile, doc *goquery.Document) {
	p.Related = getPersons(doc.Find(".insights-browse-map"))
}

func getPersons(s *goquery.Selection) []person {
	result := make([]person, 0)
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		pe := person{}
		pe.Name = s.Text()
		if pe.Name == "" {
			return
		}

		pe.Link, _ = s.Attr("href")

		result = append(result, pe)
	})

	return result
}

func getTimeFromString(date string) (time.Time, error) {
	const longForm = "January 2006"
	return time.Parse(longForm, date)
}

func getURLFromRedirect(href string) (string, error) {
	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return u.Query().Get("url"), nil
}
