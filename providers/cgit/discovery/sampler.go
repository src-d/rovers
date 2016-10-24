package discovery

import (
	"math/rand"
	"time"
)

// sampler bring the possibility of get random pages from a certain API.
// Doing this, we can eventually get all the results of that API in a certain amount of time,
// if we have limited requests per day (quotas), and the amount of pages are variable
// over the time (like a google search).
type sampler struct {
	firstIndex        int
	lastKnownEndIndex int
	multiplier        int
	r                 *rand.Rand
}

func newSampler(firstIndex int, lastKnownEndIndex int, multiplier int) *sampler {
	if multiplier < 0 {
		panic("Mutiplier must be positive")
	}

	return &sampler{
		firstIndex:        firstIndex,
		lastKnownEndIndex: lastKnownEndIndex,
		multiplier:        multiplier,
		r:                 rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// This method will return a random valid pages between firstIndex value and LastKnownEndIndex value
func (bs *sampler) RandomSampling(maxNumberOfSamples int) []int {
	pageIndexes := bs.generatePages()
	pagesCount := len(pageIndexes)
	if pagesCount <= maxNumberOfSamples {
		return pageIndexes
	}

	result := []int{}
	resultIndexes := bs.r.Perm(maxNumberOfSamples)
	for _, i := range resultIndexes {
		result = append(result, pageIndexes[i])
	}

	return result
}

func (bs *sampler) generatePages() []int {
	lastIndex := bs.findLastIndex()
	result := []int{}
	for i := bs.firstIndex; i <= lastIndex; i = i + bs.multiplier {
		result = append(result, i)
	}

	return result
}

func (bs *sampler) findLastIndex() int {
	return bs.lastKnownEndIndex
}
