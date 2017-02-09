package models

import (
	"github.com/src-d/go-kallax"
)

type URL struct {
	kallax.Model      `table:"cgit_urls"`
	kallax.Timestamps `kallax:",inline"`

	CgitUrl string
}
