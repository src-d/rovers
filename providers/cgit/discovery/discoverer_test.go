package discovery

import (
	"github.com/src-d/rovers/core"
	. "gopkg.in/check.v1"
)

type DiscovererSuite struct {
	discoverer Discoverer
}

var _ = Suite(&DiscovererSuite{
	discoverer: NewDiscoverer(core.Config.Google.SearchKey, core.Config.Google.SearchCx),
})

func (s *DiscovererSuite) TestDiscoverer_Samples(c *C) {
	samples := s.discoverer.Discover()
	c.Assert(len(samples), Equals, 100)
}
