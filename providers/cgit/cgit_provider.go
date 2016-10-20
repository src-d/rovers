package cgit

import (
	"io"
	"sync"

	"github.com/src-d/rovers/core"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
)

const (
	cgitProviderName = "cgit"
	cgitUrlField     = "CgitUrl"
	repoField        = "RepoUrls"
)

type cgitRepo struct {
	CgitUrl string
	RepoUrl string
}

type provider struct {
	cgitCollection      *mgo.Collection
	scrapers            []*scraper
	currentScraperIndex int
	mutex               *sync.Mutex
	lastRepo            string
}

func NewProvider(cgitUrls []string) *provider {
	scrapers := initializeScrapers(cgitUrls)
	p := &provider{
		cgitCollection:      initializeCollection(),
		scrapers:            scrapers,
		mutex:               &sync.Mutex{},
	}

	return p
}

func initializeScrapers(cgitUrls []string) []*scraper {
	scrapers := []*scraper{}
	for _, cgitUrl := range cgitUrls {
		scrapers = append(scrapers, newScraper(cgitUrl))
	}

	return scrapers
}

func initializeCollection() *mgo.Collection {
	cgitColl := core.NewClient().Collection(cgitProviderName)
	index := mgo.Index{
		Key: []string{"$text:" + cgitUrlField, "$text:" + repoField},
	}
	cgitColl.EnsureIndex(index)

	return cgitColl
}

func (cp *provider) setCheckpoint(cgitUrl string, repoUrl string) error {
	log15.Debug("Adding new checkpoint url", "cgitUrl", cgitUrl, "repoUrl", repoUrl)

	return cp.cgitCollection.Insert(&cgitRepo{CgitUrl: cgitUrl, RepoUrl: repoUrl})
}

func (cp *provider) alreadyProcessed(cgitUrl string, repoUrl string) (bool, error) {
	c, err := cp.cgitCollection.
		Find(&cgitRepo{CgitUrl: cgitUrl, RepoUrl: repoUrl}).Count()

	return c > 0, err
}

func (cp *provider) Next() (string, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if cp.lastRepo != "" {
		log15.Warn("Some error happens when try to call Ack(), returning the last repository again",
			"repo", cp.lastRepo)

		return cp.lastRepo, nil
	}

	for {
		currentScraper := cp.scrapers[cp.currentScraperIndex]
		cgitUrl := currentScraper.CgitUrl
		url, err := currentScraper.Next()
		switch {
		case err == io.EOF:
			cp.currentScraperIndex++
			if len(cp.scrapers) <= cp.currentScraperIndex {
				log15.Debug("All cgitUrls processed, ending provider iterator.",
					"current index", cp.currentScraperIndex)
				cp.currentScraperIndex = 0
				return "", io.EOF
			}
		case err != nil:
			return "", err
		case err == nil:
			processed, err := cp.alreadyProcessed(cgitUrl, url)
			if err != nil {
				return "", err
			}

			if processed {
				log15.Debug("Repository already processed", "cgitUrl", cgitUrl, "url", url)
			} else {
				cp.lastRepo = url
				return url, nil
			}
		}
	}
}

func (cp *provider) Ack(err error) error {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if err == nil {
		err = cp.setCheckpoint(cp.scrapers[cp.currentScraperIndex].CgitUrl, cp.lastRepo)
		if err != nil {
			return err
		} else {
			cp.lastRepo = ""
		}
	}

	return nil
}

func (cp *provider) Close() error {
	return nil
}

func (cp *provider) Name() string {
	return cgitProviderName
}
