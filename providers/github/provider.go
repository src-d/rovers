package github

import (
	"fmt"
	"io"
	"sync"

	"github.com/src-d/rovers/core"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
	"srcd.works/core.v0/models"
	"srcd.works/domain.v6/container"
	"srcd.works/domain.v6/models/social"
)

const (
	providerName         = "github"
	repositoryCollection = "repositories"

	idField       = "github_id"
	fullnameField = "fullname"
	htmlurlField  = "htmlurl"
	forkField     = "fork"

	textIndexFormat = "$text:%s"
)

type provider struct {
	repositoriesColl *mgo.Collection
	apiClient        *client
	repoStore        *social.GithubRepositoryStore
	repoCache        []*Repository
	checkpoint       int
	applyAck         func()
	mutex            *sync.Mutex
}

type Config struct {
	GithubToken string
	Database    string
}

func NewProvider(config *Config) core.RepoProvider {
	repoStore := container.GetDomainModelsSocialGithubRepositoryStore()

	return &provider{
		repositoriesColl: initRepositoriesCollection(config.Database),
		apiClient:        newClient(config.GithubToken),
		repoStore:        repoStore,
		mutex:            &sync.Mutex{},
	}
}

func initRepositoriesCollection(database string) *mgo.Collection {
	githubColl := core.NewClient(database).Collection(repositoryCollection)
	index := mgo.Index{
		Key: []string{
			fmt.Sprintf(textIndexFormat, fullnameField),
			fmt.Sprintf(textIndexFormat, htmlurlField),
			idField,
			forkField,
		},
	}
	githubColl.EnsureIndex(index)

	return githubColl
}

func (gp *provider) Name() string {
	return providerName
}

func (gp *provider) Next() (*models.Mention, error) {
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

func (*provider) repositoryRaw(repoUrl string, isFork bool) *models.Mention {
	return &models.Mention{
		Provider: providerName,
		Endpoint: repoUrl,
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

func (gp *provider) requestNextPage(since int) ([]*Repository, error) {
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
	result := Repository{}
	err := gp.repositoriesColl.Find(nil).Sort("-_id").One(&result)
	fmt.Println("RESULT:", result)
	if err == mgo.ErrNotFound {
		return 0, nil
	}

	return result.ID, err
}

func (gp *provider) saveRepos(repositories []*Repository) error {
	bulkOp := gp.repositoriesColl.Bulk()
	for _, repo := range repositories {
		bulkOp.Insert(repo)
	}
	_, err := bulkOp.Run()

	return err
}
