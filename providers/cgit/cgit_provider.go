package cgit

import (
	"io"
	"sync"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit/discovery"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	cgitProviderName = "cgit"
	cgitUrlField     = "cgiturl"
	repoField        = "repourl"
)

type cgitRepo struct {
	CgitUrl string
	RepoUrl string
	Html    string
}

type provider struct {
	cgitCollection      *mgo.Collection
	scrapers            []*scraper
	discoverer          discovery.Discoverer
	currentScraperIndex int
	mutex               *sync.Mutex
	lastRepo            *cgitRepoData
}

func NewProvider(googleKey string, googleCx string) *provider {
	p := &provider{
		cgitCollection: initializeCollection(),
		discoverer:     discovery.NewDiscoverer(googleKey, googleCx),
		scrapers:       []*scraper{},
		mutex:          &sync.Mutex{},
		lastRepo:       nil,
	}

	return p
}

func initializeCollection() *mgo.Collection {
	cgitColl := core.NewClient(cgitProviderName).Collection(cgitProviderName)
	index := mgo.Index{
		Key: []string{"$text:" + cgitUrlField, "$text:" + repoField},
	}
	cgitColl.EnsureIndex(index)

	return cgitColl
}

func (cp *provider) setCheckpoint(cgitUrl string, repo *cgitRepoData) error {
	log15.Debug("Adding new checkpoint url", "cgitUrl", cgitUrl, "repoUrl", repo.RepoUrl)

	return cp.cgitCollection.Insert(&cgitRepo{CgitUrl: cgitUrl, RepoUrl: repo.RepoUrl, Html: repo.Html})
}

func (cp *provider) alreadyProcessed(cgitUrl string, repo *cgitRepoData) (bool, error) {
	c, err := cp.cgitCollection.Find(bson.M{cgitUrlField: cgitUrl, repoField: repo.RepoUrl}).Count()

	return c > 0, err
}

func (cp *provider) getAllCgitUrlsAlreadyProcessed() ([]string, error) {
	cgitUrls := []string{}
	err := cp.cgitCollection.Find(nil).Distinct(cgitUrlField, &cgitUrls)

	return cgitUrls, err
}

func (cp *provider) fillScrapersIfNecessary() error {
	if len(cp.scrapers) == 0 {
		cgitUrlsSet := map[string]struct{}{}

		alreadyProcessedCgitUrls, err := cp.getAllCgitUrlsAlreadyProcessed()
		if err != nil {
			return err
		}

		cgitUrls := cp.discoverer.Samples()
		cp.addToSet(cgitUrlsSet, cgitUrls)
		cp.addToSet(cgitUrlsSet, alreadyProcessedCgitUrls)
		for u := range cgitUrlsSet {
			log15.Info("Adding new Scraper", "cgit url", u)
			cp.scrapers = append(cp.scrapers, newScraper(u))
		}

		if len(cp.scrapers) == 0 {
			// No scrapers found, sending an EOF because we have no data
			return io.EOF
		}
	}

	return nil
}

func (cp *provider) addToSet(set map[string]struct{}, slice []string) {
	for _, e := range slice {
		_, ok := set[e]
		if !ok {
			set[e] = struct{}{}
		}
	}
}

func (cp *provider) Next() (string, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if cp.lastRepo != nil {
		log15.Warn("Some error happens when try to call Ack(), returning the last repository again",
			"repo", cp.lastRepo.RepoUrl)

		return cp.lastRepo.RepoUrl, nil
	}

	err := cp.fillScrapersIfNecessary()
	if err != nil {
		return "", err
	}

	for {
		currentScraper := cp.scrapers[cp.currentScraperIndex]
		cgitUrl := currentScraper.CgitUrl
		repoData, err := currentScraper.Next()
		switch {
		case err == io.EOF:
			cp.currentScraperIndex++
			if len(cp.scrapers) <= cp.currentScraperIndex {
				log15.Debug("All cgitUrls processed, ending provider iterator.",
					"current index", cp.currentScraperIndex)
				cp.scrapers = []*scraper{}
				cp.currentScraperIndex = 0
				return "", io.EOF
			}
		case err != nil:
			return "", err
		case err == nil:
			processed, err := cp.alreadyProcessed(cgitUrl, repoData)
			if err != nil {
				return "", err
			}

			if processed {
				log15.Debug("Repository already processed", "cgitUrl", cgitUrl, "url", repoData)
			} else {
				cp.lastRepo = repoData
				return repoData.RepoUrl, nil
			}
		}
	}
}

func (cp *provider) Ack(err error) error {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if err == nil && cp.lastRepo != nil {
		err = cp.setCheckpoint(cp.scrapers[cp.currentScraperIndex].CgitUrl, cp.lastRepo)
		if err != nil {
			return err
		} else {
			cp.lastRepo = nil
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
