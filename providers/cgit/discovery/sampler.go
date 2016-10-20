package discovery

import "math/rand"

type sampler struct {
	firstIndex        int
	lastKnownEndIndex int
	multiplier        int
}

func newSampler(firstIndex int, lastKnownEndIndex int, multiplier int) *sampler {
	return &sampler{
		firstIndex:        firstIndex,
		lastKnownEndIndex: lastKnownEndIndex,
		multiplier:        multiplier,
	}
}

func (bs *sampler) RandomSampling(maxNumberOfSamplings int) []int {
	pageIndexes := bs.generatePages()
	pagesCount := len(pageIndexes)
	if pagesCount <= maxNumberOfSamplings {
		return pageIndexes
	}

	result := []int{}
	resultIndexes := rand.Perm(maxNumberOfSamplings)
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
	// TODO not necessary by now boolean search. Change this in the future.
	return bs.lastKnownEndIndex
}
