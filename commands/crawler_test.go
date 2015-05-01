package commands

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CrawlerSuite struct{}

var _ = Suite(&CrawlerSuite{})

func (s *CrawlerSuite) TestCrawler_SearchGithub(c *C) {
	cr := NewCrawler()
	r, err := cr.SearchGithub("Máximo Cuadros")
	c.Assert(err, IsNil)
	c.Assert(r.FullName, Equals, "Máximo Cuadros")
}
