package readers

import (
	"sort"

	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestAugurReader_SearchByEmail(c *C) {
	c.Skip("not used")
	a := NewAugurInsightsAPI(client.NewClient(true))
	r, res, err := a.SearchByEmail("nawar.alsafar126@gmail.com")

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(r.LastStatus, Equals, 200)
	sort.Strings(r.Name)
	c.Assert(r.Name, DeepEquals, []string{"Nawar Alsafar", "Noir Alsafar"})
}

func (s *SourcesSuite) TestAugurReader_SearchByEmail_BadRequest(c *C) {
	c.Skip("not used")
	a := NewAugurInsightsAPI(client.NewClient(true))
	r, res, err := a.SearchByEmail("foo")

	c.Assert(r, IsNil)
	c.Assert(err, Equals, ErrUnexpectedStatusCode)
	c.Assert(res.StatusCode, Equals, 400)
}
