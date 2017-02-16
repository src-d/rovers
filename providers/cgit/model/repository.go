package model

import "github.com/src-d/go-kallax"

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
