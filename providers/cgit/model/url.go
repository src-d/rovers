package model

import (
	"gopkg.in/src-d/go-kallax.v1"
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
