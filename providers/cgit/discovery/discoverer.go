package discovery

import "gopkg.in/inconshreveable/log15.v2"

const (
	firstIndex      = 1
	lastIndex       = 91
	elementsPerPage = 10

	maxNumberOfPages = 100
)

type Discoverer interface {
	Samples() []string
}

type DefaultDiscoverer struct {
	sampler   *sampler
	googleApi *googleCseApi
}

func NewDiscoverer(googleKey string, googleCx string) Discoverer {
	return &DefaultDiscoverer{
		sampler:   newSampler(firstIndex, lastIndex, elementsPerPage),
		googleApi: newGoogleCseApi(googleKey, googleCx),
	}
}

func (d *DefaultDiscoverer) Samples() []string {
	samples := d.sampler.RandomSampling(maxNumberOfPages)
	cgitUrls := []string{}
	for _, s := range samples {
		p, err := d.googleApi.GetPage(s)
		if err != nil {
			log15.Warn("Error obtaining google page", "Error", err)
			continue
		}

		for _, i := range p.Items {
			cgitUrls = append(cgitUrls, i.Link)
		}
	}

	return cgitUrls
}
