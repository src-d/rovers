package readers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/tyba/oss/domain/models/social"
	"github.com/tyba/oss/domain/models/social/linkedin"
	"github.com/tyba/oss/sources/social/http"

	"github.com/PuerkitoBio/goquery"
)

type LinkedInReader struct {
	client *http.Client
}

func NewLinkedInReader(client *http.Client) *LinkedInReader {
	return &LinkedInReader{client}
}

func (l *LinkedInReader) GetProfileByURL(url string) (*social.LinkedInProfile, error) {
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

	p := &social.LinkedInProfile{URL: url, Created: time.Now()}
	l.fillConnections(p, doc)
	l.fillBasicInfo(p, doc)
	l.fillExperience(p, doc)
	l.fillLanguages(p, doc)
	l.fillSkills(p, doc)
	l.fillEducation(p, doc)
	l.fillPatents(p, doc)
	l.fillProjects(p, doc)
	l.fillPublications(p, doc)
	l.fillAlsoViewed(p, doc)

	return p, nil
}

func (l *LinkedInReader) fillConnections(p *social.LinkedInProfile, doc *goquery.Document) {
	c := doc.Find(".profile-overview .member-connections strong").Text()
	if c == "500+" {
		p.Connections = 500
	} else {
		p.Connections, _ = strconv.Atoi(c)
	}
}

func (l *LinkedInReader) fillBasicInfo(p *social.LinkedInProfile, doc *goquery.Document) {
	p.FullName = doc.Find(".full-name").Text()
	p.Title = doc.Find(".title").Text()
	p.Industry = doc.Find("#location .industry").Text()
	p.Locality = doc.Find("#location .locality").Text()

	p.Websites = make(map[string]string, 0)
	doc.Find("#overview-summary-websites").Find("li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		url, _ := l.getURLFromRedirect(href)

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

func (l *LinkedInReader) fillExperience(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Experience = make([]linkedin.Experience, 0)
	doc.Find("#background-experience").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		e := linkedin.Experience{}
		e.Title = s.Find("h4").Text()
		e.Company = s.Find("h5").Text()
		e.Link, _ = s.Find("h5").Find("a").Attr("href")
		e.Locality = s.Find(".locality").Text()
		e.Summary = s.Find(".description").Text()

		s.Find("time").Each(func(i int, s *goquery.Selection) {
			date, _ := l.getTimeFromString(s.Text())
			if i == 0 {
				e.StartDate = date
			} else {
				e.EndDate = date
			}
		})

		p.Experience = append(p.Experience, e)
	})
}

func (l *LinkedInReader) fillLanguages(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Languages = make(map[string]string, 0)
	doc.Find("#languages").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		language := s.Find("h4").Text()
		level := s.Find(".languages-proficiency").Text()

		p.Languages[language] = level
	})
}

func (l *LinkedInReader) fillSkills(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Skills = make([]string, 0)
	doc.Find(".endorse-item ").Each(func(i int, s *goquery.Selection) {
		if _, ok := s.Attr("id"); !ok {
			p.Skills = append(p.Skills, s.Text())
		}
	})
}

func (l *LinkedInReader) fillEducation(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Education = make([]linkedin.Education, 0)
	doc.Find("#background-education").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		e := linkedin.Education{}
		e.Summary = s.Find("h4").Text()
		e.Degree = s.Find(".degree").Text()
		e.Major = s.Find(".major").Text()
		e.Date = s.Find(".education-date").Text()
		e.Extra = s.Find("h5").Text()

		p.Education = append(p.Education, e)
	})
}

func (l *LinkedInReader) fillPatents(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Patents = make([]linkedin.Patent, 0)
	doc.Find("#background-patents").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		pa := linkedin.Patent{}
		pa.Title = s.Find("h4 span:first-of-type").Text()
		pa.Id = s.Find("h5").Text()
		pa.Date = s.Find(".patents-date").Text()
		pa.Summary = s.Find(".description").Text()
		pa.Persons = l.getPersons(s.Find(".associated-endorsements"))

		p.Patents = append(p.Patents, pa)
	})
}

func (l *LinkedInReader) fillProjects(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Projects = make([]linkedin.Project, 0)
	doc.Find("#background-projects").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		pr := linkedin.Project{}
		pr.Title = s.Find("h4 span:first-of-type").Text()
		href, _ := s.Find("h4").Find("a").Attr("href")
		pr.Link, _ = l.getURLFromRedirect(href)
		pr.Date = s.Find("time").Text()
		pr.Summary = s.Find(".description").Text()
		pr.Persons = l.getPersons(s.Find(".associated-endorsements"))

		p.Projects = append(p.Projects, pr)
	})
}

func (l *LinkedInReader) fillPublications(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Publications = make([]linkedin.Publication, 0)
	doc.Find("#background-publications").Find(".section-item").Each(func(i int, s *goquery.Selection) {
		pu := linkedin.Publication{}
		pu.Title = s.Find("h4").Text()
		href, _ := s.Find("h4").Find("a").Attr("href")
		pu.Link, _ = l.getURLFromRedirect(href)
		pu.Date = s.Find(".publication-date").Text()
		pu.Summary = s.Find(".description").Text()
		pu.Persons = l.getPersons(s.Find(".associated-endorsements"))

		p.Publications = append(p.Publications, pu)
	})
}

func (l *LinkedInReader) fillAlsoViewed(p *social.LinkedInProfile, doc *goquery.Document) {
	p.Related = l.getPersons(doc.Find(".insights-browse-map"))
}

func (l *LinkedInReader) getPersons(s *goquery.Selection) []linkedin.Person {
	result := make([]linkedin.Person, 0)
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		pe := linkedin.Person{}
		pe.Name = s.Text()
		if pe.Name == "" {
			return
		}

		pe.Link, _ = s.Attr("href")

		result = append(result, pe)
	})

	return result
}

func (l *LinkedInReader) getTimeFromString(date string) (time.Time, error) {
	const longForm = "January 2006"
	return time.Parse(longForm, date)
}

func (l *LinkedInReader) getURLFromRedirect(href string) (string, error) {
	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return u.Query().Get("url"), nil
}
