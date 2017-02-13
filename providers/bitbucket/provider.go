package bitbucket

import (
	"database/sql"
	"sync"
	"time"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers"
	"github.com/src-d/rovers/providers/bitbucket/model"

	"github.com/src-d/go-kallax"
	"gopkg.in/inconshreveable/log15.v2"
	coreModels "srcd.works/core.v0/model"
)

const (
	providerName = "bitbucket"

	gitScm        = "git"
	httpsCloneKey = "https"

	firstCheckpoint    = ""
	minRequestDuration = time.Hour / 1000
)

type provider struct {
	repositoryStore   *model.RepositoryStore
	client            *client

	mutex             *sync.Mutex
	repositoriesCache model.Repositories
	lastCheckpoint    string
	applyAck          func()
}

func NewProvider(database *sql.DB) core.RepoProvider {
	return &provider{
		repositoryStore: model.NewRepositoryStore(database),
		client:          newClient(),
		mutex:           &sync.Mutex{},
		lastCheckpoint:  firstCheckpoint,
	}
}

func (p *provider) isInit() bool {
	return p.repositoriesCache == nil && p.lastCheckpoint == firstCheckpoint
}

func (p *provider) needsMoreData() bool {
	return len(p.repositoriesCache) == 0
}

func (p *provider) repositoryRaw(r *model.Repository) *coreModels.Mention {
	aliases := []string{}
	mainRepository := ""
	for _, c := range r.Links.Clone {
		if c.Name == httpsCloneKey {
			mainRepository = c.Href
		}
		aliases = append(aliases, c.Href)
	}
	if mainRepository == "" {
		log15.Error("no https repositories found", "clone urls", r.Links.Clone)
	}

	return &coreModels.Mention{
		Endpoint: mainRepository,
		Provider: providerName,
		VCS:      coreModels.GIT,
		Context: providers.ContextBuilder{}.
			Fork(r.Parent != nil).
			Aliases(aliases),
	}
}

func (p *provider) initializeCheckpoint() error {
	result, err := p.repositoryStore.FindOne(
		model.NewRepositoryQuery().
			Order(kallax.Asc(model.Schema.Repository.CreatedAt)),
	)

	switch err {
	case nil:
		log15.Info("checkpoint found", "checkpoint", result.Next)
		p.lastCheckpoint = result.Next
	case kallax.ErrNotFound:
		p.lastCheckpoint = firstCheckpoint
	default:
		return err
	}

	return nil
}

func (p *provider) Next() (*coreModels.Mention, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for {
		if p.isInit() {
			err := p.initializeCheckpoint()
			if err != nil {
				return nil, err
			}
		}

		if p.needsMoreData() {
			err := p.requestNextPage()
			if err != nil {
				return nil, err
			}
		}

		r, repositories := p.repositoriesCache[0], p.repositoriesCache[1:]
		if r.Scm == gitScm {
			p.applyAck = func() {
				p.repositoriesCache = repositories
			}

			return p.repositoryRaw(r), nil
		} else {
			log15.Debug("non git repository found", "repository", r.FullName, "scm", r.Scm)
			p.repositoriesCache = repositories
		}
	}

}

func (p *provider) requestNextPage() error {
	start := time.Now()
	defer func() {
		needsWait := minRequestDuration - time.Since(start)
		if needsWait > 0 {
			log15.Debug("waiting", "duration", needsWait)
			time.Sleep(needsWait)
		}
	}()
	response, err := p.client.Repositories(p.lastCheckpoint)
	if err != nil {
		return err
	}
	p.lastCheckpoint = response.Next
	p.repositoriesCache = response.Repositories

	return p.saveRepositories(response)
}

func (p *provider) saveRepositories(resp *response) error {
	return p.repositoryStore.Transaction(func(store *model.RepositoryStore) error {
		// TODO implements bulk operations in kallax
		for _, repo := range resp.Repositories {
			repo.Next = resp.Next
			if _, err := store.Save(repo); err != nil {
				return err
			}
		}

		return nil
	})
}

func (p *provider) Ack(err error) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if err == nil {
		if p.applyAck != nil {
			p.applyAck()
		}
	} else {
		log15.Warn("error when watcher tried to send last url. Not applying ack", "error", err)
	}

	return nil
}

func (p *provider) Close() error {
	return nil
}

func (p *provider) Name() string {
	return providerName
}
