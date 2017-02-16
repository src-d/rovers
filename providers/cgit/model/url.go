package model

import (
	"github.com/src-d/go-kallax"
)

func newURL() *URL {
	return &URL{ID: kallax.NewULID()}
}

type URL struct {
	ID                kallax.ULID `pk:""`
	kallax.Model      `table:"cgit_urls"`
	kallax.Timestamps `kallax:",inline"`

	CgitUrl string
}
