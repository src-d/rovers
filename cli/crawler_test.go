package cli

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CrawlerSuite struct{}

var _ = Suite(&CrawlerSuite{})

func (s *CrawlerSuite) Test_normalize(c *C) {
	cr := NewCrawler()

	r := cr.normalize("Fóo. a. bar-qúx")
	c.Assert(r, Equals, "foo bar-qux")
}

func (s *CrawlerSuite) TestCrawler_SearchGithub(c *C) {
	cr := NewCrawler()
	r, err := cr.SearchGithub("Máximo Cuadros")
	c.Assert(err, IsNil)
	c.Assert(r.FullName, Equals, "Máximo Cuadros")
}
