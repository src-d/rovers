package discovery

import (
	"math/rand"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type SamplerSuite struct {
	sampler *sampler
}

var _ = Suite(&SamplerSuite{})

func (s *SamplerSuite) SetUpTest(c *C) {
	s.sampler = newSampler(1, 91, 10)
}
func (s *SamplerSuite) TestSampler_findLastIndex(c *C) {
	lastIndex := s.sampler.findLastIndex()
	c.Assert(lastIndex, Equals, 91)
}
func (s *SamplerSuite) TestSampler_RandomSampling(c *C) {
	samples := s.sampler.RandomSampling(10)
	c.Assert(len(samples), Equals, 10)
}

func (s *SamplerSuite) TestSampler_RandomSampling_Random(c *C) {
	s.sampler = &sampler{
		firstIndex:        1,
		lastKnownEndIndex: 91,
		multiplier:        10,
		r:                 rand.New(rand.NewSource(42)),
	}
	samplesOne := s.sampler.RandomSampling(5)
	expectedSamplesOne := []int{1, 11, 31, 41, 21}
	c.Assert(samplesOne, DeepEquals, expectedSamplesOne)
	c.Assert(len(samplesOne), Equals, 5)

	samplesTwo := s.sampler.RandomSampling(5)
	expectedSamplesTwo := []int{31, 11, 21, 41, 1}
	c.Assert(samplesTwo, DeepEquals, expectedSamplesTwo)
	c.Assert(len(samplesTwo), Equals, 5)
}

func (s *SamplerSuite) TestSampler_RandomSampling_TryToGetMoreSamplings(c *C) {
	samples := s.sampler.RandomSampling(30)
	c.Assert(len(samples), Equals, 10)
}

func (s *SamplerSuite) TestSampler_RandomSampling_OnlyFiveSamplings(c *C) {
	samples := s.sampler.RandomSampling(5)
	c.Assert(len(samples), Equals, 5)
}

func (s *SamplerSuite) TestSampler_RandomSampling_LastKnownIndexSmallerThanFirstIndex(c *C) {
	s.sampler = newSampler(100, 0, 10)

	samples := s.sampler.RandomSampling(11)
	c.Assert(len(samples), Equals, 0)
}

func (s *SamplerSuite) TestSampler_RandomSampling_NegativeMultiplier(c *C) {

	defer func() {
		if r := recover(); r == nil {
			c.Errorf("The code did not panic")
		}
	}()

	s.sampler = newSampler(0, 100, -10)
}

func (s *SamplerSuite) TestSampler_RandomSampling_DifferentConstructor(c *C) {
	s.sampler = newSampler(0, 100, 10)
	samples := s.sampler.RandomSampling(11)
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
