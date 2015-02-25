package sources

import (
	"strconv"
	"strings"
	"time"

	"github.com/tyba/opensource-search/sources/social/http"

	"github.com/PuerkitoBio/goquery"
)

type Github struct {
	client *http.Client
}

func NewGithub(client *http.Client) *Github {
	return &Github{client}
}

func (g *Github) GetProfileByURL(url string) (*GithubProfile, error) {
	req, err := http.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, _, err := g.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	return NewGithubProfile(url, doc), nil
}

type GithubProfile struct {
	Created            time.Time
	Url                string
	Username           string
	FullName           string
	Location           string
	Email              string
	Web                string
	WorksFor           string
	JoinDate           time.Time
	Description        string
	Organizations      []string
	Repositories       []*repository
	Contributions      []*repository
	Followers          int
	Starred            int
	Following          int
	TotalContributions int
}

type repository struct {
	Owner       string
	Name        string
	Description string
	Url         string
	Stars       int
}

func NewGithubProfile(url string, doc *goquery.Document) *GithubProfile {
	g := &GithubProfile{Url: url, Created: time.Now()}
	g.fillBasicInfo(doc)
	g.fillOrganizations(doc)
	g.fillRepositories(doc)
	g.fillContributions(doc)
	g.fillStats(doc)

	return g
}

func (g *GithubProfile) fillBasicInfo(doc *goquery.Document) {
	g.Username = doc.Find(".vcard-username").Text()
	g.FullName = doc.Find(".vcard-fullname").Text()
	g.Location = doc.Find("[itemprop='homeLocation']").Text()
	g.Email = doc.Find(".email").Text()
	g.Web = doc.Find("[itemprop='url']").Text()
	g.WorksFor = doc.Find("[itemprop='worksFor']").Text()
	g.Description, _ = doc.Find("[property='og:description']").Attr("content")

	date, _ := doc.Find(".join-date").Attr("datetime")
	g.JoinDate, _ = time.Parse("2006-01-02T15:04:05Z", date)
}

func (g *GithubProfile) fillOrganizations(doc *goquery.Document) {
	g.Organizations = make([]string, 0)
	doc.Find("[itemprop='follows']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		g.Organizations = append(g.Organizations, href)
	})
}

func (g *GithubProfile) fillRepositories(doc *goquery.Document) {
	g.Repositories = make([]*repository, 0)
	doc.Find(".one-half:first-of-type").Find(".mini-repo-list-item").Each(func(i int, s *goquery.Selection) {
		r := &repository{}
		r.Owner = g.Username
		r.Name = s.Find(".repo").Text()
		r.Description = s.Find(".repo-description").Text()
		r.Url, _ = s.Attr("href")
		r.Stars, _ = strconv.Atoi(strings.Trim(s.Find(".stars").Text(), " \n"))

		g.Repositories = append(g.Repositories, r)
	})
}

func (g *GithubProfile) fillContributions(doc *goquery.Document) {
	g.Contributions = make([]*repository, 0)
	doc.Find(".one-half:last-of-type").Find(".mini-repo-list-item").Each(func(i int, s *goquery.Selection) {
		r := &repository{}
		r.Owner = s.Find(".owner").Text()
		r.Name = s.Find(".repo").Text()
		r.Description = s.Find(".repo-description").Text()
		r.Url, _ = s.Attr("href")
		r.Stars, _ = strconv.Atoi(strings.Trim(s.Find(".stars").Text(), " \n"))

		g.Contributions = append(g.Contributions, r)
	})
}

func (g *GithubProfile) fillStats(doc *goquery.Document) {
	g.Followers, _ = strconv.Atoi(doc.Find(".vcard-stat-count").Eq(0).Text())
	g.Starred, _ = strconv.Atoi(doc.Find(".vcard-stat-count").Eq(1).Text())
	g.Following, _ = strconv.Atoi(doc.Find(".vcard-stat-count").Eq(2).Text())

	contribs := doc.Find(".contrib-number").Text()
	contribs = strings.Replace(contribs, ",", "", -1)
	total := strings.Index(contribs, " ")
	if len(contribs) != 0 {
		g.TotalContributions, _ = strconv.Atoi(contribs[:total])
	}
}
