package model

import "github.com/src-d/go-kallax"

type Repository struct {
	kallax.Model      `table:"cgit"`
	kallax.Timestamps `kallax:",inline"`

	CgitURL string
	URL     string
	Aliases []string
	HTML    string
}
