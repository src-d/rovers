package cgit

import (
	"database/sql"
	"errors"
	"io"
	goURL "net/url"
	"sync"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit/model"

	. "gopkg.in/check.v1"
	rcore "gopkg.in/src-d/core-retrieval.v0"
	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
)

type CgitProviderSuite struct {
	DB *sql.DB
}

var _ = Suite(&CgitProviderSuite{})

func (s *CgitProviderSuite) SetUpTest(c *C) {
	DB := rcore.Database()
	s.DB = DB

	err := core.DropTables(DB, providerName, "cgit_urls")
	c.Assert(err, IsNil)

	err = core.CreateCgitTables(DB)
	c.Assert(err, IsNil)
}

func (s *CgitProviderSuite) newProvider(cgitUrls ...string) *provider {
	return &provider{
		repositoriesStore: model.NewRepositoryStore(s.DB),
		urlsStore:         model.NewURLStore(s.DB),
		searcher:          &dummySearcher{cgitUrls},
		backoff:           getBackoff(),
		scrapers:          []*scraper{},
		mutex:             &sync.Mutex{},
	}
}

func (s *CgitProviderSuite) TestCgitProvider_WhenFinishScraping(c *C) {
	provider := s.newProvider("https://a3nm.net/git/")

	var err error
	var url *rmodel.Mention
	count := 0
	for err == nil {
		url, err = provider.Next()
		if err == nil {
			ackErr := provider.Ack(nil)
			c.Assert(ackErr, IsNil)
		}
		count++
	}

	c.Assert(count, Not(Equals), 0)
	c.Assert(url, IsNil)
	c.Assert(err, Equals, io.EOF)

}

func (s *CgitProviderSuite) TestCgitProvider_WhenAckIsError(c *C) {
	provider := s.newProvider("https://a3nm.net/git/")

	urlOne, err := provider.Next()
	ackErr := provider.Ack(errors.New("OOPS"))
	c.Assert(err, IsNil)
	c.Assert(ackErr, IsNil)

	urlTwo, err := provider.Next()
	ackErr = provider.Ack(nil)
	c.Assert(err, IsNil)
	c.Assert(ackErr, IsNil)

	urlTree, err := provider.Next()
	c.Assert(err, IsNil)

	c.Assert(urlOne, DeepEquals, urlTwo)
	c.Assert(urlTwo, Not(DeepEquals), urlTree)
}

func (s *CgitProviderSuite) TestCgitProvider_NotSendAlreadySended(c *C) {
	provider := s.newProvider("https://a3nm.net/git/")

	urlOne, err := provider.Next()
	ackErr := provider.Ack(nil)
	c.Assert(err, IsNil)
	c.Assert(ackErr, IsNil)

	provider = s.newProvider("https://a3nm.net/git/")

	urlTwo, err := provider.Next()
	ackErr = provider.Ack(nil)
	c.Assert(err, IsNil)
	c.Assert(ackErr, IsNil)

	c.Assert(urlOne, Not(DeepEquals), urlTwo)
}

func (s *CgitProviderSuite) TestCgitProvider_IterateAllUrls(c *C) {
	provider := s.newProvider("https://a3nm.net/git/", "https://ongardie.net/git/")
	maxIndex := 0
	for {
		_, err := provider.Next()
		if provider.currentScraperIndex > maxIndex {
			maxIndex = provider.currentScraperIndex
		}
		if err == io.EOF {
			break
		}
		c.Assert(err, IsNil)
		ackErr := provider.Ack(nil)
		c.Assert(ackErr, IsNil)
	}
	c.Assert(maxIndex, Equals, 1)
	c.Assert(provider.currentScraperIndex, Equals, 0)
	c.Assert(len(provider.scrapers), Equals, 0)
}

func (s *CgitProviderSuite) TestCgitProvider_ScrapersWithDifferentUrls(c *C) {
	provider := s.newProvider("https://a3nm.net/git/", "https://a3nm.net/git/", "https://ongardie.net/git/")
	_, err := provider.Next()
	c.Assert(err, IsNil)
	c.Assert(len(provider.scrapers), Equals, 2)
}

func (s *CgitProviderSuite) TestCgitProvider_Retries(c *C) {
	provider := s.newProvider()
	provider.scrapers = []*scraper{newScraper("https://badurl.com")}
	_, err := provider.Next()
	c.Assert(err, NotNil)
	c.Assert(provider.backoff.Attempt(), Equals, float64(1))
	_, err = provider.Next()
	c.Assert(err, NotNil)
	c.Assert(provider.backoff.Attempt(), Equals, float64(2))
}

func (s *CgitProviderSuite) TestCgitProvider_RetriesBadUrl(c *C) {
	provider := s.newProvider("https://badurl.com")
	_, err := provider.Next()
	c.Assert(err, Equals, io.EOF)
	c.Assert(provider.backoff.Attempt(), Equals, float64(0))
}

func (s *CgitProviderSuite) TestCgitProvider_CgitUrlsNotDuplicated(c *C) {
	provider := s.newProvider("https://a3nm.net/git/", "https://a3nm.net/git/", "https://ongardie.net/git/")
	_, err := provider.Next()
	c.Assert(err, IsNil)

	uCount, err := provider.urlsStore.Count(model.NewURLQuery())
	c.Assert(err, IsNil)
	c.Assert(uCount, Equals, int64(2))

	provider = s.newProvider("https://a3nm.net/git/")
	_, err = provider.Next()
	c.Assert(err, IsNil)
	uCount, err = provider.urlsStore.Count(model.NewURLQuery())
	c.Assert(err, IsNil)
	c.Assert(uCount, Equals, int64(2))

	provider = s.newProvider("http://pkgs.fedoraproject.org/cgit/rpms/")
	_, err = provider.Next()
	c.Assert(err, IsNil)
	uCount, err = provider.urlsStore.Count(model.NewURLQuery())
	c.Assert(err, IsNil)
	c.Assert(uCount, Equals, int64(3))

	provider = s.newProvider()
	_, err = provider.Next()
	c.Assert(err, IsNil)
	c.Assert(len(provider.scrapers), Equals, 3)
}

type dummySearcher struct {
	urls []string
}

func (d *dummySearcher) Search(query string) ([]*goURL.URL, error) {
	result := []*goURL.URL{}
	for _, s := range d.urls {
		u, _ := goURL.Parse(s)
		result = append(result, u)
	}
	return result, nil
}
