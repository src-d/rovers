package github

import (
	"fmt"

	"github.com/src-d/rovers/core"
	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
)

const (
	httpsUrl = "https://github.com/%s.git"
	sshUrl   = "git@github.com:%s.git"
	gitUrl   = "git://github.com/%s"
)

func getMention(repoName string, isFork bool) *rmodel.Mention {
	gu := fmt.Sprintf(gitUrl, repoName)
	su := fmt.Sprintf(sshUrl, repoName)
	hu := fmt.Sprintf(httpsUrl, repoName)

	return &rmodel.Mention{
		Provider: core.GithubProviderName,
		Endpoint: gu,
		VCS:      rmodel.GIT,
		IsFork:   &isFork,
		Aliases:  []string{gu, su, hu},
	}
}
