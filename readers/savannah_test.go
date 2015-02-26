package readers

import (
	"github.com/tyba/opensource-search/sources/social/http"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestSavannah_GetRepositories(c *C) {
	a := NewSavannahReader(http.NewClient(true))
	r, err := a.GetRepositories()
	c.Assert(err, IsNil)
	c.Assert(len(r.Results) > 0, Equals, true)
}
