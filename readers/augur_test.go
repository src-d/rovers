package readers

import (
	"github.com/tyba/srcd-rovers/client"

	. "gopkg.in/check.v1"
)

type SourcesSuite struct{}

var _ = Suite(&SourcesSuite{})

func (s *SourcesSuite) TestAugurReader_SearchByEmail(c *C) {
	a := NewAugurInsightsAPI(client.NewClient(true))
	r, res, err := a.SearchByEmail("nawar.alsafar126@gmail.com")

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(r.Status, Equals, 200)
}

func (s *SourcesSuite) TestAugurReader_SearchByEmailBadRequest(c *C) {
	a := NewAugurInsightsAPI(client.NewClient(true))
	r, res, err := a.SearchByEmail("foo")

	c.Assert(r, IsNil)
	c.Assert(err, Equals, ErrUnexpectedStatusCode)
	c.Assert(res.StatusCode, Equals, 422)
}
