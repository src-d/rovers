package websearch

import "net/url"

type Searcher interface {
	Search(query string) ([]*url.URL, error)
}
