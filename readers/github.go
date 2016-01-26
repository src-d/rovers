package readers

import (
	"strconv"
	"strings"
	"time"

	"github.com/src-d/rovers/client"
	"gop.kg/src-d/domain@v3/models/social"
	"gop.kg/src-d/domain@v3/models/social/github"

	"github.com/PuerkitoBio/goquery"
)

type GithubWebCrawler struct {
	client *client.Client
}

func NewGithubWebCrawler(client *client.Client) *GithubWebCrawler {
	return &GithubWebCrawler{client}
}

func (g *GithubWebCrawler) GetProfileByURL(url string) (*social.GithubProfile, error) {
	req, err := client.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, res, err := g.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, client.NotFound
	}

	p := &social.GithubProfile{Url: url, Created: time.Now()}
	g.fillOrganizationInfo(doc, p)
	if p.Organization == true {
		g.fillMembers(doc, p)
		return p, nil
	}

	g.fillBasicInfo(doc, p)
	g.fillOrganizations(doc, p)
	g.fillRepositories(doc, p)
	g.fillContributions(doc, p)
	g.fillStats(doc, p)

	return p, nil
}

func (g *GithubWebCrawler) fillOrganizationInfo(doc *goquery.Document, p *social.GithubProfile) {
	urlParts := strings.Split(p.Url, "/")
	if len(urlParts) >= 4 {
		p.Username = urlParts[3]
	}
	p.FullName = doc.Find(".org-name span").Text()
	p.Location = doc.Find("[itemprop='location']").Text()
	p.Email = doc.Find("[itemprop='email']").Text()
	p.Web = doc.Find("[itemprop='url']").Text()
	if p.FullName != "" {
		p.Organization = true
	}
}

func (g *GithubWebCrawler) fillMembers(doc *goquery.Document, p *social.GithubProfile) {
	p.Members = make([]string, 0)
	doc.Find(".member-avatar-group .avatar").Each(func(i int, s *goquery.Selection) {
		username, _ := s.Attr("alt")
		if len(username) != 0 && username[0] == '@' {
			p.Members = append(p.Members, username[1:])
		}
	})
}

func (g *GithubWebCrawler) fillBasicInfo(doc *goquery.Document, p *social.GithubProfile) {
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

func (g *GithubWebCrawler) fillOrganizations(doc *goquery.Document, p *social.GithubProfile) {
	p.Organizations = make([]string, 0)
	doc.Find("[itemprop='follows']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		p.Organizations = append(p.Organizations, href)
	})
}

func (g *GithubWebCrawler) fillRepositories(doc *goquery.Document, p *social.GithubProfile) {
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

func (g *GithubWebCrawler) fillContributions(doc *goquery.Document, p *social.GithubProfile) {
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

func (g *GithubWebCrawler) fillStats(doc *goquery.Document, p *social.GithubProfile) {
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
