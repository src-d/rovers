package readers

import (
	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestTwitter_GetProfileByHandle(c *C) {
	a := NewTwitterReader(client.NewClient(true))
	g, err := a.GetProfileByURL("https://twitter.com/mcuadros_")
	c.Assert(err, IsNil)
	c.Assert(g.URL, Equals, "https://twitter.com/mcuadros_")
	c.Assert(g.Handle, Equals, "mcuadros_")
	c.Assert(g.FullName, Equals, "Máximo Cuadros")
	c.Assert(g.Bio, Equals, "(╯°□°）╯︵ ┻━┻")
	c.Assert(g.Location, Equals, "Madrid, Spain")
	c.Assert(g.Web, Equals, "http://github.com/mcuadros")

	c.Assert(g.Tweets > 0, Equals, true)
	c.Assert(g.Following > 0, Equals, true)
	c.Assert(g.Followers > 0, Equals, true)
	c.Assert(g.Favorites > 0, Equals, true)
}
