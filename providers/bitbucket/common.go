package bitbucket

import (
	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/bitbucket/model"
	"gopkg.in/inconshreveable/log15.v2"

	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
)

const (
	gitScm        = "git"
	httpsCloneKey = "https"
)

func getMention(r *model.Repository) *rmodel.Mention {
	aliases := []string{}
	mainRepository := ""
	for _, c := range r.Links.Clone {
		if c.Name == httpsCloneKey {
			mainRepository = c.Href
		}
		aliases = append(aliases, c.Href)
	}

	if mainRepository == "" {
		if len(aliases) > 0 {
			mainRepository = aliases[0]
		} else {
			log15.Error("no https repositories found", "clone urls", r.Links.Clone)
		}
	}

	isFork := r.Parent != nil && r.Parent.UUID != ""
	return &rmodel.Mention{
		Endpoint: mainRepository,
		Provider: core.BitbucketProviderName,
		VCS:      rmodel.GIT,
		IsFork:   &isFork,
		Aliases:  aliases,
	}
}
