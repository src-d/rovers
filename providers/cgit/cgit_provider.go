package cgit

import (
	"io"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/sourcegraph/go-vcsurl"
	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit/discovery"
	"gop.kg/src-d/domain@v6/container"
	"gop.kg/src-d/domain@v6/models/repository"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	cgitProviderName     = "cgit"
	repositoryCollection = "repositories"

	cgitURLField    = "cgiturl"
	repositoryField = "repourl"

	maxDurationToRetry = 16 * time.Second
	minDurationToRetry = 1 * time.Second
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
	backoff             *backoff.Backoff
	currentScraperIndex int
	mutex               *sync.Mutex
	lastRepo            *cgitRepoData
}

func getBackoff() *backoff.Backoff {
	return &backoff.Backoff{
		Jitter: true,
		Factor: 2,
		Max:    maxDurationToRetry,
		Min:    minDurationToRetry,
	}
}

func NewProvider(googleKey string, googleCx string) *provider {
	p := &provider{
		cgitCollection: initializeCollection(),
		discoverer:     discovery.NewDiscoverer(googleKey, googleCx),
		backoff:        getBackoff(),
		scrapers:       []*scraper{},
		mutex:          &sync.Mutex{},
		lastRepo:       nil,
	}

	return p
}

func initializeCollection() *mgo.Collection {
	cgitColl := core.NewClient(container.Config.MongoDb.Database.Cgit).Collection(repositoryCollection)
	index := mgo.Index{
		Key: []string{"$text:" + cgitURLField, "$text:" + repositoryField},
	}
	cgitColl.EnsureIndex(index)

	return cgitColl
}

func (cp *provider) setCheckpoint(cgitUrl string, repo *cgitRepoData) error {
	log15.Debug("Adding new checkpoint url", "cgitUrl", cgitUrl, "repoUrl", repo.RepoUrl)

	return cp.cgitCollection.Insert(&cgitRepo{CgitUrl: cgitUrl, RepoUrl: repo.RepoUrl, Html: repo.Html})
}

func (cp *provider) alreadyProcessed(cgitUrl string, repo *cgitRepoData) (bool, error) {
	c, err := cp.cgitCollection.Find(bson.M{cgitURLField: cgitUrl, repositoryField: repo.RepoUrl}).Count()

	return c > 0, err
}

func (cp *provider) getAllCgitUrlsAlreadyProcessed() ([]string, error) {
	cgitUrls := []string{}
	err := cp.cgitCollection.Find(nil).Distinct(cgitURLField, &cgitUrls)

	return cgitUrls, err
}

func (cp *provider) fillScrapers() {
	cgitUrlsSet := map[string]struct{}{}
	alreadyProcessedCgitUrls, err := cp.getAllCgitUrlsAlreadyProcessed()
	if err != nil {
		log15.Error("Error getting cgit urls from database", "error", err)
	}

	cp.discoverer.Reset()
	cgitUrls := cp.discoverer.Discover()
	cp.joinUnique(cgitUrlsSet, cgitUrls, alreadyProcessedCgitUrls)
	for u := range cgitUrlsSet {
		log15.Info("Adding new Scraper", "cgit url", u)
		cp.scrapers = append(cp.scrapers, newScraper(u))
	}
}

func (cp *provider) joinUnique(set map[string]struct{}, slices ...[]string) {
	for _, slice := range slices {
		for _, e := range slice {
			set[e] = struct{}{}
		}
	}
}

func (cp *provider) Next() (*repository.Raw, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if cp.lastRepo != nil {
		log15.Warn("Some error happens when try to call Ack(), returning the last repository again",
			"repo", cp.lastRepo.RepoUrl)

		return cp.repositoryRaw(cp.lastRepo.RepoUrl), nil
	}

	if cp.isFirst() {
		cp.fillScrapers()
		if len(cp.scrapers) == 0 {
			log15.Warn("No scrapers found, sending an EOF because we have no data")
			return nil, io.EOF
		}
	}

	for {
		currentScraper := cp.scrapers[cp.currentScraperIndex]
		cgitUrl := currentScraper.CgitUrl
		repoData, err := currentScraper.Next()
		switch {
		case err == io.EOF:
			cp.nextScraper()
			if len(cp.scrapers) <= cp.currentScraperIndex {
				log15.Debug("All cgitUrls processed, ending provider iterator.",
					"current index", cp.currentScraperIndex)
				cp.reset()
				return nil, io.EOF
			}
		case err != nil:
			log15.Error("error on scraper.next", "cgitUrl", currentScraper.CgitUrl, "error", err)
			cp.handleRetries()
			return nil, err
		case err == nil:
			cp.backoff.Reset()
			processed, err := cp.alreadyProcessed(cgitUrl, repoData)
			if err != nil {
				return nil, err
			}

			if processed {
				log15.Debug("Repository already processed", "cgitUrl", cgitUrl, "url", repoData.RepoUrl)
			} else {
				cp.lastRepo = repoData
				return cp.repositoryRaw(repoData.RepoUrl), nil
			}
		}
	}
}

func (cp *provider) handleRetries() {
	tts := cp.backoff.Duration()
	log15.Info("Sleeping before next scraper request",
		"retries", cp.backoff.Attempt(), "time to sleep", tts)
	time.Sleep(tts)

	if tts >= maxDurationToRetry {
		log15.Warn("Scraper request failed too many times. Skipping to the next scraper",
			"retries", cp.backoff.Attempt())
		cp.nextScraper()
	}
}

func (*provider) repositoryRaw(repoUrl string) *repository.Raw {
	return &repository.Raw{
		Status:   repository.Initial,
		Provider: cgitProviderName,
		URL:      repoUrl,
		VCS:      vcsurl.Git,
	}
}

func (cp *provider) nextScraper() {
	cp.backoff.Reset()
	cp.currentScraperIndex++
}

func (cp *provider) isFirst() bool {
	return len(cp.scrapers) == 0
}

func (cp *provider) reset() {
	cp.scrapers = []*scraper{}
	cp.currentScraperIndex = 0
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
