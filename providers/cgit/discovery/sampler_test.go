package discovery

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type SamplerSuite struct {
	boolSearch *sampler
}

var _ = Suite(&SamplerSuite{})

func (s *SamplerSuite) SetUpTest(c *C) {
	s.boolSearch = newSampler(1, 91, 10)
}
func (s *SamplerSuite) TestSampler_findLastIndex(c *C) {
	lastIndex := s.boolSearch.findLastIndex()
	c.Assert(lastIndex, Equals, 91)
}
func (s *SamplerSuite) TestSampler_RandomSampling(c *C) {
	samples := s.boolSearch.RandomSampling(10)
	c.Assert(len(samples), Equals, 10)
}

func (s *SamplerSuite) TestSampler_RandomSampling_Random(c *C) {
	samplesOne := s.boolSearch.RandomSampling(5)
	samplesTwo := s.boolSearch.RandomSampling(5)
	c.Assert(len(samplesOne), Equals, 5)
	c.Assert(len(samplesTwo), Equals, 5)
	c.Assert(samplesOne, Not(DeepEquals), samplesTwo)
}

func (s *SamplerSuite) TestSampler_RandomSampling_TryToGetMoreSamplings(c *C) {
	samples := s.boolSearch.RandomSampling(30)
	c.Assert(len(samples), Equals, 10)
}

func (s *SamplerSuite) TestSampler_RandomSampling_OnlyFiveSamplings(c *C) {
	samples := s.boolSearch.RandomSampling(5)
	c.Assert(len(samples), Equals, 5)
}

func (s *SamplerSuite) TestSampler_RandomSampling_DifferentConstructor(c *C) {
	s.boolSearch = newSampler(0, 100, 10)
	samples := s.boolSearch.RandomSampling(11)
	c.Assert(len(samples), Equals, 11)

	cont := contains(samples, 20)
	c.Assert(cont, Equals, true)

	cont = contains(samples, 100)
	c.Assert(cont, Equals, true)

	cont = contains(samples, 101)
	c.Assert(cont, Equals, false)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
