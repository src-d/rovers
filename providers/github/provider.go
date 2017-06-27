package github

import (
	"database/sql"
	"io"
	"sync"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers"
	"github.com/src-d/rovers/providers/github/model"

	"gopkg.in/inconshreveable/log15.v2"
	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
	"gopkg.in/src-d/go-kallax.v1"
)

const (
	providerName = "github"
)

type provider struct {
	repositoriesStore *model.RepositoryStore
	apiClient         *client
	repoCache         []*model.Repository
	checkpoint        int
	applyAck          func()
	mutex             *sync.Mutex
}

func NewProvider(githubToken string, DB *sql.DB) core.RepoProvider {
	return &provider{
		repositoriesStore: model.NewRepositoryStore(DB),
		apiClient:         newClient(githubToken),
		mutex:             &sync.Mutex{},
	}
}

func (gp *provider) Name() string {
	return providerName
}

func (gp *provider) Next() (*rmodel.Mention, error) {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	switch len(gp.repoCache) {
	case 0:
		if gp.checkpoint == 0 {
			log15.Info("checkpoint empty, trying to get checkpoint")
			c, err := gp.getLastRepoId()
			if err != nil {
				log15.Error("error getting checkpoint from database", "error", err)
				return nil, err
			}
			gp.checkpoint = c
		}
		log15.Info("no repositories into cache, getting more repositories", "checkpoint", gp.checkpoint)
		repos, err := gp.requestNextPage(gp.checkpoint)
		if err != nil {
			log15.Error("something bad happens getting more repositories", "error", err)
			return nil, err
		}
		if len(repos) != 0 {
			gp.repoCache = repos
		} else {
			log15.Info("no more repositories, sending EOF")
			return nil, io.EOF
		}
	}

	x, repoCache := gp.repoCache[0], gp.repoCache[1:]
	gp.applyAck = func() {
		gp.repoCache = repoCache
	}

	return gp.repositoryRaw(x.HTMLURL+".git", x.Fork), nil
}

func (*provider) repositoryRaw(repoUrl string, isFork bool) *rmodel.Mention {
	return &rmodel.Mention{
		Provider: providerName,
		Endpoint: repoUrl,
		VCS:      rmodel.GIT,
		Context:  providers.ContextBuilder{}.Fork(isFork),
	}
}

func (gp *provider) Ack(err error) error {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	if err == nil {
		if gp.applyAck != nil {
			gp.applyAck()
		}
	} else {
		log15.Warn("error when watcher tried to send last url. Not applying ack", "error", err)
	}

	return nil
}

func (gp *provider) Close() error {
	return nil
}

func (gp *provider) requestNextPage(since int) ([]*model.Repository, error) {
	resp, err := gp.apiClient.Repositories(since)
	if err != nil {
		return nil, err
	}

	gp.checkpoint = resp.Next

	if err := gp.saveRepos(resp.Repositories); err != nil {
		return nil, err
	}

	return resp.Repositories, nil
}

func (gp *provider) getLastRepoId() (int, error) {
	result, err := gp.repositoriesStore.FindOne(model.NewRepositoryQuery().
		Order(kallax.Desc(model.Schema.Repository.CreatedAt)))

	if err == kallax.ErrNotFound {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return result.GithubID, nil
}

func (gp *provider) saveRepos(repositories []*model.Repository) error {
	return gp.repositoriesStore.Transaction(func(s *model.RepositoryStore) error {
		for _, repo := range repositories {
			repo.ID = kallax.NewULID()
			err := s.Insert(repo)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
