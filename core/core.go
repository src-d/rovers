package core

import "gopkg.in/src-d/core-retrieval.v0/model"

const (
	GithubProviderName    = "github"
	CgitProviderName      = "cgit"
	BitbucketProviderName = "bitbucket"
)

type RepoProvider interface {
	Next() (*model.Mention, error)
	Ack(error) error
	Close() error
	Name() string
}
