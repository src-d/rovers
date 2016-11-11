package bitbucket

import (
	"sync"
	"time"

	"github.com/src-d/rovers/core"

	"github.com/sourcegraph/go-vcsurl"
	"gop.kg/src-d/domain@v6/models/repository"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
)

const (
	providerName = "bitbucket"

	repositoriesCollection = "repositories"

	idDescKey = "-_id"

	scmField      = "scm"
	fullNameField = "fullname"

	gitScm        = "git"
	httpsCloneKey = "https"

	firstCheckpoint    = ""
	minRequestDuration = time.Hour / 1000
)

type provider struct {
	repoCollection *mgo.Collection
	client         *client

	mutex             *sync.Mutex
	repositoriesCache repositories
	lastCheckpoint    string
	applyAck          func()
}

func NewProvider(database string) core.RepoProvider {
	return &provider{
		repoCollection: initializeRepositoriesCollection(database),
		client:         newClient(),
		mutex:          &sync.Mutex{},
		lastCheckpoint: firstCheckpoint,
	}
}

func initializeRepositoriesCollection(database string) *mgo.Collection {
	coll := core.NewClient(database).Collection(repositoriesCollection)
	coll.EnsureIndexKey(scmField, fullNameField)

	return coll
}

func (p *provider) isInit() bool {
	return p.repositoriesCache == nil && p.lastCheckpoint == firstCheckpoint
}

func (p *provider) needsMoreData() bool {
	return len(p.repositoriesCache) == 0
}

func (p *provider) repositoryRaw(r *bitbucketRepository) *repository.Raw {
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

	return &repository.Raw{
		Status:   repository.Initial,
		Provider: providerName,
		URL:      mainRepository,
		IsFork:   r.Parent != nil,
		VCS:      vcsurl.Git,
		Aliases:  aliases,
	}
}

func (p *provider) initializeCheckpoint() error {
	result := bitbucketRepository{}
	err := p.repoCollection.Find(nil).Sort(idDescKey).One(&result)

	switch err {
	case nil:
		log15.Info("checkpoint found", "checkpoint", result.Next)
		p.lastCheckpoint = result.Next
	case mgo.ErrNotFound:
		p.lastCheckpoint = firstCheckpoint
	default:
		return err
	}

	return nil
}

func (p *provider) Next() (*repository.Raw, error) {
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
	bulkOp := p.repoCollection.Bulk()
	for _, repo := range resp.Repositories {
		repo.Next = resp.Next
		bulkOp.Insert(repo)
	}
	_, err := bulkOp.Run()

	return err
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
