package discovery

import (
	"os"

	. "gopkg.in/check.v1"
)

type DiscovererSuite struct {
	discoverer Discoverer
}

var _ = Suite(&DiscovererSuite{
	discoverer: NewDiscoverer(os.Getenv(envKey), os.Getenv(envCx)),
})

func (s *DiscovererSuite) TestDiscoverer_Samples(c *C) {
	samples := s.discoverer.Discover()
	c.Assert(len(samples), Equals, 100)
}
