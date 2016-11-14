package cgit

import (
	"io"

	. "gopkg.in/check.v1"
)

const (
	gitUrl   = "git://pkgs.fedoraproject.org/rpms/0ad.git"
	sshUrl   = "ssh://pkgs.fedoraproject.org/rpms/0ad.git"
	httpUrl  = "http://pkgs.fedoraproject.org/git/rpms/0ad.git"
	httpsUrl = "https://pkgs.fedoraproject.org/git/rpms/0ad.git"
	otherUrl = "other://pkgs.fedoraproject.org/git/rpms/0ad.git"
	noResult = ""
)

type CgitScraperSuite struct {
}

var _ = Suite(&CgitScraperSuite{})

func (s *CgitScraperSuite) TestCgitScraper_Next_CorrectMainPage(c *C) {
	scraper := newScraper("http://pkgs.fedoraproject.org/cgit/")
	u, err := scraper.Next()
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)
	c.Assert(u.Html, Not(Equals), "")
	c.Assert(u.RepositoryURL, Not(Equals), "")
}

func (s *CgitScraperSuite) TestCgitScraper_Next_CorrectMainPageWithNoPages(c *C) {
	scraper := newScraper("http://git.mate-desktop.org/")
	u, err := scraper.Next()
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)
}

func (s *CgitScraperSuite) TestCgitScraper_Next_IncorrectMainPage(c *C) {
	scraper := newScraper("http://git.mate-desktop.org/libmateweather/")
	u, err := scraper.Next()
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)
}

func (s *CgitScraperSuite) TestCgitScraper_Next_IncorrectPage(c *C) {
	scraper := newScraper("https://golang.org/ref/spec")
	u, err := scraper.Next()
	c.Assert(err, NotNil)
	c.Assert(u, IsNil)
}

func (s *CgitScraperSuite) TestCgitScraper_Next_EOF(c *C) {
	scraper := newScraper("https://a3nm.net/git/")
	var err error = nil
	count := 0
	for err != io.EOF {
		_, err = scraper.Next()
		count++
	}
	c.Assert(count, Not(Equals), 0)
	c.Assert(err, Equals, io.EOF)
	u, err := scraper.Next()
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)
}

func (s *CgitScraperSuite) TestCgitScraper_repoPageWithNoRepos(c *C) {
	scraper := newScraper("https://a3nm.net/git/")
	var err error = nil
	count := 0
	for err != io.EOF {
		_, err = scraper.Next()
		count++
	}
	c.Assert(count, Not(Equals), 0)
	c.Assert(err, Equals, io.EOF)
	u, err := scraper.Next()
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)
}

func (s *CgitScraperSuite) TestCgitScraper_repoPageWithNoHttpRepos(c *C) {
	scraper := newScraper("http://cgit.openembedded.org/")
	url, err := scraper.Next()
	c.Assert(url, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *CgitScraperSuite) TestCgitScraper_mainPage(c *C) {
	scraper := newScraper("")

	urlTests := []*inOutCase{
		{in: []string{sshUrl, gitUrl, httpUrl, httpsUrl, otherUrl}, out: httpsUrl},
		{in: []string{otherUrl}, out: noResult},
		{in: []string{httpUrl}, out: httpUrl},
		{in: []string{gitUrl, httpUrl}, out: gitUrl},
		{in: nil, out: noResult},
	}

	for _, d := range urlTests {
		c.Assert(scraper.mainUrl(d.in), Equals, d.out)
	}
}

type inOutCase struct {
	in  []string
	out string
}
