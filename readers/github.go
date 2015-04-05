package readers

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tyba/opensource-search/domain/models/social"
	"github.com/tyba/opensource-search/domain/models/social/github"
	"github.com/tyba/opensource-search/sources/social/http"
)

type GithubReader struct {
	client *http.Client
}

func NewGithubReader(client *http.Client) *GithubReader {
	return &GithubReader{client}
}

func (g *GithubReader) GetProfileByURL(url string) (*social.GithubProfile, error) {
	req, err := http.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, _, err := g.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	p := &social.GithubProfile{Url: url, Created: time.Now()}
	g.fillBasicInfo(doc, p)
	g.fillOrganizations(doc, p)
	g.fillRepositories(doc, p)
	g.fillContributions(doc, p)
	g.fillStats(doc, p)

	return p, nil
}

func (g *GithubReader) fillBasicInfo(doc *goquery.Document, p *social.GithubProfile) {
	p.Username = doc.Find(".vcard-username").Text()
	p.FullName = doc.Find(".vcard-fullname").Text()
	p.Location = doc.Find("[itemprop='homeLocation']").Text()
	p.Email = doc.Find(".email").Text()
	p.Web = doc.Find("[itemprop='url']").Text()
	p.WorksFor = doc.Find("[itemprop='worksFor']").Text()
	p.Description, _ = doc.Find("[property='og:description']").Attr("content")

	date, _ := doc.Find(".join-date").Attr("datetime")
	p.JoinDate, _ = time.Parse("2006-01-02T15:04:05Z", date)
}

func (g *GithubReader) fillOrganizations(doc *goquery.Document, p *social.GithubProfile) {
	p.Organizations = make([]string, 0)
	doc.Find("[itemprop='follows']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		p.Organizations = append(p.Organizations, href)
	})
}

func (g *GithubReader) fillRepositories(doc *goquery.Document, p *social.GithubProfile) {
	p.Repositories = make([]github.Repository, 0)
	doc.Find(".one-half:first-of-type").Find(".mini-repo-list-item").Each(func(i int, s *goquery.Selection) {
		r := github.Repository{}
		r.Owner = p.Username
		r.Name = s.Find(".repo").Text()
		r.Description = s.Find(".repo-description").Text()
		r.Url, _ = s.Attr("href")
		r.Stars, _ = strconv.Atoi(strings.Trim(s.Find(".stars").Text(), " \n"))

		p.Repositories = append(p.Repositories, r)
	})
}

func (g *GithubReader) fillContributions(doc *goquery.Document, p *social.GithubProfile) {
	p.Contributions = make([]github.Repository, 0)
	doc.Find(".one-half:last-of-type").Find(".mini-repo-list-item").Each(func(i int, s *goquery.Selection) {
		r := github.Repository{}
		r.Owner = s.Find(".owner").Text()
		r.Name = s.Find(".repo").Text()
		r.Description = s.Find(".repo-description").Text()
		r.Url, _ = s.Attr("href")
		r.Stars, _ = strconv.Atoi(strings.Trim(s.Find(".stars").Text(), " \n"))

		p.Contributions = append(p.Contributions, r)
	})
}

func (g *GithubReader) fillStats(doc *goquery.Document, p *social.GithubProfile) {
	p.Followers, _ = strconv.Atoi(doc.Find(".vcard-stat-count").Eq(0).Text())
	p.Starred, _ = strconv.Atoi(doc.Find(".vcard-stat-count").Eq(1).Text())
	p.Following, _ = strconv.Atoi(doc.Find(".vcard-stat-count").Eq(2).Text())

	contribs := doc.Find(".contrib-number").Text()
	contribs = strings.Replace(contribs, ",", "", -1)
	total := strings.Index(contribs, " ")
	if len(contribs) != 0 {
		p.TotalContributions, _ = strconv.Atoi(contribs[:total])
	}
}
