package cgit

import (
	"errors"
	"fmt"
	"io"
	goURL "net/url"
	"strings"
	"time"

	"github.com/src-d/rovers/utils"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	repoHttpUrlSelector = "table.list tr a[href^=http]"
	paginationSelector  = "ul.pager li a"
	pagesUrlSelector    = "div.content table tr td.toplevel-repo a, td.sublevel-repo a"
	mainPageSelector    = "td.logo a"

	hrefAttr = "href"

	httpTimeout = 30 * time.Second

	prefixGit   = "git://"
	prefixHttps = "https://"
	prefixHttp  = "http://"
	prefixSsh   = "ssh://"
)

var gitPrefixesByPreference = []string{prefixHttps, prefixGit, prefixHttp, prefixSsh}

type page struct {
	RepositoryURL string
	Aliases       []string
	Html          string
}

type scraper struct {
	URL             string
	firstIteration  bool
	pageUrls        []string
	repositoryPages []string
	goqueryClient   *utils.GoqueryClient
}

func newScraper(cgitUrl string) *scraper {
	return &scraper{
		URL:             cgitUrl,
		firstIteration:  true,
		pageUrls:        []string{},
		repositoryPages: []string{},
		goqueryClient:   utils.NewDefaultGoqueryClient(),
	}
}

func (cs *scraper) Next() (*page, error) {
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

func (cs *scraper) repoFound(repo *page) bool {
	return repo.RepositoryURL != ""
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
	log15.Debug("first execution, adding more page URLs", "cgit URL", cs.URL)
	mainPage, err := mainPage(cs.URL, cs.goqueryClient)
	if err != nil {
		return err
	}
	pageUrls, err := cs.paginationUrls(mainPage)
	if err != nil {
		return err
	}
	if len(pageUrls) == 0 {
		log15.Debug("main page with no pagination. Scraping main page directly", "cgit URL", cs.URL)
		cs.pageUrls = []string{mainPage}
	} else {
		cs.pageUrls = pageUrls
	}
	cs.firstIteration = false

	return nil
}

func (cs *scraper) refreshPages() error {
	log15.Debug("repository pages are empty, adding more.", "cgit URL", cs.URL)
	pageUrl, pageUrls := cs.pageUrls[0], cs.pageUrls[1:]
	repoPageUrls, err := cs.repoPageUrls(pageUrl)
	if err != nil {
		return err
	}
	cs.pageUrls = pageUrls
	cs.repositoryPages = repoPageUrls

	return nil
}

func (cs *scraper) getRepo() (*page, error) {
	if len(cs.repositoryPages) == 0 {
		return nil, errors.New("no repository pages found")
	}
	repoPage, repositoryPages := cs.repositoryPages[0], cs.repositoryPages[1:]
	repoData, err := cs.repo(repoPage)
	if err != nil {
		return nil, err
	}
	cs.repositoryPages = repositoryPages

	return repoData, nil
}

func getAllMainCgitUrls(possibleCgitURLs []string) []string {
	goqueryClient := utils.NewDefaultGoqueryClient()
	cgitUrls := []string{}
	for _, pcu := range possibleCgitURLs {
		cgitUrl, err := mainPage(pcu, goqueryClient)
		if err != nil {
			log15.Warn("Error getting cgit main page url", "error", err)
			continue
		}
		cgitUrls = append(cgitUrls, cgitUrl)
	}
	return cgitUrls
}

func mainPage(cgitUrl string, gqClient *utils.GoqueryClient) (string, error) {
	mainDoc, err := gqClient.NewDocument(cgitUrl)
	if err != nil {
		return "", err
	}

	href, exists := mainDoc.Find(mainPageSelector).Attr(hrefAttr)
	if !exists {
		return "", fmt.Errorf("tried to scrape a non correct cgit URL: %v", cgitUrl)
	}
	urlType, err := utils.BaseURL(cgitUrl)
	if err != nil {
		return "", err
	}
	mainUrl := &goURL.URL{
		Scheme: urlType.Scheme,
		Host:   urlType.Host,
		Path:   href,
	}
	if mainUrl.String() != cgitUrl {
		log15.Info("we are not in the main page, getting data from main page", "input URL", cgitUrl, "main URL", mainUrl)
		return mainPage(mainUrl.String(), gqClient)
	} else {
		return mainUrl.String(), nil
	}
}

func (cs *scraper) scrapeMain(initUrl string, selector string,
	fun func(s *goquery.Selection, baseUrl string) string) ([]string, error) {
	urlsToScrape := []string{}
	baseUrl, err := utils.BaseURL(cs.URL)
	if err != nil {
		return nil, err
	}
	document, err := cs.goqueryClient.NewDocument(initUrl)
	if err != nil {
		return nil, err
	}
	document.Find(selector).Each(
		func(i int, selection *goquery.Selection) {
			u := fun(selection, baseUrl.String())
			if u != "" {
				urlsToScrape = append(urlsToScrape, u)
			}
		})

	return urlsToScrape, nil
}

func (cs *scraper) paginationUrls(mainPageUrl string) ([]string, error) {
	return cs.scrapeMain(mainPageUrl, paginationSelector,
		func(s *goquery.Selection, baseUrl string) string {
			pageUrl, exists := s.Attr(hrefAttr)
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
			repoPageUrlPath, exists := s.Attr(hrefAttr)
			return baseUrl + repoPageUrlPath
			if exists {
				return baseUrl + repoPageUrlPath
			} else {
				return ""
			}
		})
}

func (cs *scraper) repo(repoUrl string) (*page, error) {
	document, err := cs.goqueryClient.NewDocument(repoUrl)
	if err != nil {
		return nil, err
	}

	html, err := document.Html()
	if err != nil {
		return nil, err
	}

	urls := document.Find(repoHttpUrlSelector).Map(func(i int, s *goquery.Selection) string {
		// all '<a>' html tags have href in this case
		href, _ := s.Attr(hrefAttr)
		return href
	})

	return &page{
		RepositoryURL: cs.mainUrl(urls),
		Aliases:       urls,
		Html:          html,
	}, nil
}

func (cs *scraper) mainUrl(urls []string) string {
	for _, prefix := range gitPrefixesByPreference {
		for _, u := range urls {
			if strings.HasPrefix(u, prefix) {
				return u
			}
		}
	}

	return ""
}

func (cs *scraper) newDocument(url string) (*goquery.Document, error) {
	resp, err := cs.goqueryClient.Get(url)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(resp)
}
