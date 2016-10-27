package cgit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	repoHttpUrlSelector = "table.list tr a[href^=http]"
	paginationSelector  = "ul.pager li a"
	pagesUrlSelector    = "div.content table tr td.toplevel-repo a, td.sublevel-repo a"
	mainPageSelector    = "td.logo a"

	httpTimeout = 30 * time.Second
)

type cgitRepoData struct {
	RepoUrl string
	Html    string
}

type scraper struct {
	CgitUrl         string
	firstIteration  bool
	pageUrls        []string
	repositoryPages []string
	client          http.Client
}

func newScraper(cgitUrl string) *scraper {
	return &scraper{
		CgitUrl:         cgitUrl,
		firstIteration:  true,
		pageUrls:        []string{},
		repositoryPages: []string{},
		client: http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (cs *scraper) Next() (*cgitRepoData, error) {
	for {
		if cs.isStart() {
			if err := cs.initialize(); err != nil {
				return nil, err
			}
		}

		if cs.isEnd() {
			cs.firstIteration = true
			return nil, io.EOF
		}

		if cs.needMorePages() {
			if err := cs.refreshPages(); err != nil {
				return nil, err
			}
		}

		repoData, err := cs.getRepo()
		if err != nil {
			return nil, err
		}

		if cs.repoFound(repoData) {
			return repoData, nil
		}
	}
}

func (cs *scraper) repoFound(repo *cgitRepoData) bool {
	return repo.RepoUrl != ""
}

func (cs *scraper) needMorePages() bool {
	return len(cs.repositoryPages) == 0 && len(cs.pageUrls) != 0
}

func (cs *scraper) isEnd() bool {
	return len(cs.pageUrls) == 0 && len(cs.repositoryPages) == 0 && !cs.firstIteration
}

func (cs *scraper) isStart() bool {
	return len(cs.pageUrls) == 0 && cs.firstIteration
}

func (cs *scraper) initialize() error {
	log15.Debug("First execution, adding more page URLs", "cgitPage", cs.CgitUrl)
	mainPage, err := cs.mainPage(cs.CgitUrl)
	if err != nil {
		return err
	}
	pageUrls, err := cs.paginationUrls(mainPage)
	if err != nil {
		return err
	}
	if len(pageUrls) == 0 {
		log15.Debug("Main page with no pagination. Scraping main page directly", "cgitPage", cs.CgitUrl)
		cs.pageUrls = []string{mainPage}
	} else {
		cs.pageUrls = pageUrls
	}
	cs.firstIteration = false

	return nil
}

func (cs *scraper) refreshPages() error {
	log15.Debug("Repository pages are empty, adding more.", "cgitPage", cs.CgitUrl)
	pageUrl, pageUrls := cs.pageUrls[0], cs.pageUrls[1:]
	repoPageUrls, err := cs.repoPageUrls(pageUrl)
	if err != nil {
		return err
	}
	cs.pageUrls = pageUrls
	cs.repositoryPages = repoPageUrls

	return nil
}

func (cs *scraper) getRepo() (*cgitRepoData, error) {
	if len(cs.repositoryPages) == 0 {
		return nil, fmt.Errorf("No repository pages found: %v", cs.repositoryPages)
	}
	repoPage, repositoryPages := cs.repositoryPages[0], cs.repositoryPages[1:]
	repoData, err := cs.repo(repoPage)
	if err != nil {
		return nil, err
	}
	cs.repositoryPages = repositoryPages

	return repoData, nil
}

func (cs *scraper) baseUrl() (*url.URL, error) {
	urlType, err := url.Parse(cs.CgitUrl)
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme: urlType.Scheme,
		Host:   urlType.Host,
	}, nil
}

func (cs *scraper) mainPage(cgitUrl string) (string, error) {
	mainDoc, err := cs.newDocument(cgitUrl)
	if err != nil {
		return "", err
	}
	href, exists := mainDoc.Find(mainPageSelector).Attr("href")
	if !exists {
		return "", fmt.Errorf("Tried to scrape a non correct cgit url: %v", cgitUrl)
	}
	urlType, err := cs.baseUrl()
	if err != nil {
		return "", err
	}
	mainUrl := &url.URL{
		Scheme: urlType.Scheme,
		Host:   urlType.Host,
		Path:   href,
	}
	if mainUrl.String() != cgitUrl {
		log15.Info("We are not in the main page, getting data from main page", "inputUrl", cgitUrl, "mainPage", mainUrl)
		return cs.mainPage(mainUrl.String())
	} else {
		return mainUrl.String(), nil
	}
}

func (cs *scraper) scrapeMain(initUrl string, selector string,
	fun func(s *goquery.Selection, baseUrl string) string) ([]string, error) {
	urlsToScrape := []string{}
	baseUrl, err := cs.baseUrl()
	if err != nil {
		return nil, err
	}
	document, err := cs.newDocument(initUrl)
	if err != nil {
		return nil, err
	}
	document.Find(selector).Each(
		func(i int, selection *goquery.Selection) {
			url := fun(selection, baseUrl.String())
			if url != "" {
				urlsToScrape = append(urlsToScrape, url)
			}
		})

	return urlsToScrape, nil
}

func (cs *scraper) paginationUrls(mainPageUrl string) ([]string, error) {
	return cs.scrapeMain(mainPageUrl, paginationSelector,
		func(s *goquery.Selection, baseUrl string) string {
			pageUrl, exists := s.Attr("href")
			if exists {
				return baseUrl + pageUrl
			} else {
				return ""
			}
		})
}

func (cs *scraper) repoPageUrls(pageUrl string) ([]string, error) {
	return cs.scrapeMain(pageUrl, pagesUrlSelector,
		func(s *goquery.Selection, baseUrl string) string {
			repoPageUrlPath, exists := s.Attr("href")
			return baseUrl + repoPageUrlPath
			if exists {
				return baseUrl + repoPageUrlPath
			} else {
				return ""
			}
		})
}

func (cs *scraper) repo(repoUrl string) (*cgitRepoData, error) {
	document, err := cs.newDocument(repoUrl)
	if err != nil {
		return nil, err
	}

	html, err := document.Html()
	if err != nil {
		return nil, err
	}

	r := ""
	document.Find(repoHttpUrlSelector).EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			repo, exists := selection.Attr("href")
			if exists {
				if strings.HasPrefix(repo, "https") {
					r = repo
					return false
				}
				r = repo
			}
			return true
		})

	return &cgitRepoData{
		RepoUrl: r,
		Html:    html,
	}, nil
}

func (cs *scraper) newDocument(url string) (*goquery.Document, error) {
	resp, err := cs.client.Get(url)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(resp)
}
