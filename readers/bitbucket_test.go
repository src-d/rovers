package readers

import (
	"github.com/tyba/opensource-search/sources/social/http"
	"net/url"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestBitbucket_GetRepositories(c *C) {
	a := NewBitbucketReader(http.NewClient(true))
	g, err := a.GetRepositories(url.Values{})
	c.Assert(err, IsNil)
	c.Assert(g.Next.Query().Get("page"), Equals, "2")
}
