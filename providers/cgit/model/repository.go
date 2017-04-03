package model

import "gopkg.in/src-d/go-kallax.v1"

func newRepository() *Repository {
	return &Repository{ID: kallax.NewULID()}
}

type Repository struct {
	ID                kallax.ULID `pk:""`
	kallax.Model      `table:"cgit"`
	kallax.Timestamps `kallax:",inline"`

	CgitURL string
	URL     string
	Aliases []string
	HTML    string
}
