package github

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sourcegraph/go-vcsurl"
	api "github.com/src-d/go-github/github"
	"github.com/src-d/rovers/core"
	"golang.org/x/oauth2"
	"gop.kg/src-d/domain@v6/container"
	"gop.kg/src-d/domain@v6/models/repository"
	"gop.kg/src-d/domain@v6/models/social"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
)

const (
	minRequestDuration = time.Hour / 5000

	providerName         = "github"
	repositoryCollection = "repositories"

	idField       = "id"
	fullnameField = "fullname"
	htmlurlField  = "htmlurl"
	forkField     = "fork"

	textIndexFormat = "$text:%s"
)

type githubProvider struct {
	repositoriesCollection *mgo.Collection
	apiClient              *api.Client
	repoStore              *social.GithubRepositoryStore
	repoCache              []*api.Repository
	checkpoint             int
	applyAck               func()
	mutex                  *sync.Mutex
}

type GithubConfig struct {
	GithubToken string
	Database    string
}

func NewProvider(config *GithubConfig) core.RepoProvider {
	httpClient := http.DefaultClient
	if config.GithubToken != "" {
		token := &oauth2.Token{AccessToken: config.GithubToken}
		httpClient = oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))
	} else {
		log15.Warn("Creating anonymous http client. No GitHub token provided.")
	}
	apiClient := api.NewClient(httpClient)
	repoStore := container.GetDomainModelsSocialGithubRepositoryStore()

	return &githubProvider{
		initializeCollection(config.Database),
		apiClient,
		repoStore,
		[]*api.Repository{},
		0,
		nil,
		&sync.Mutex{},
	}
}

func initializeCollection(database string) *mgo.Collection {
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

func (gp *githubProvider) Name() string {
	return providerName
}

func (gp *githubProvider) Next() (*repository.Raw, error) {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	switch len(gp.repoCache) {
	case 0:
		if gp.checkpoint == 0 {
			log15.Info("Checkpoint empty, trying to get checkpoint")
			c, err := gp.getLastRepoId()
			if err != nil {
				log15.Error("Error getting checkpoint from database", "error", err)
				return nil, err
			}
			gp.checkpoint = c
		}
		log15.Info("No repositories into cache, getting more repositories", "checkpoint", gp.checkpoint)
		repos, err := gp.requestNextPage(gp.checkpoint)
		if err != nil {
			log15.Error("Something bad happens getting more repositories", "error", err)
			return nil, err
		}
		if len(repos) != 0 {
			gp.repoCache = repos
		} else {
			log15.Info("No more repos, sending EOF")
			return nil, io.EOF
		}
	}

	x, repoCache := gp.repoCache[0], gp.repoCache[1:]
	gp.applyAck = func() {
		gp.repoCache = repoCache
	}

	return gp.repositoryRaw(*x.HTMLURL+".git", *x.Fork), nil
}

func (*githubProvider) repositoryRaw(repoUrl string, isFork bool) *repository.Raw {
	return &repository.Raw{
		Status:   repository.Initial,
		Provider: providerName,
		URL:      repoUrl,
		IsFork:   isFork,
		VCS:      vcsurl.Git,
	}
}

func (gp *githubProvider) Ack(err error) error {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	if err == nil {
		if gp.applyAck != nil {
			gp.applyAck()
		}
	} else {
		log15.Warn("Error when watcher tried to send last url. Not applying ack", "error", err)
	}

	return nil
}

func (gp *githubProvider) Close() error {
	return nil
}

func (gp *githubProvider) requestNextPage(since int) ([]*api.Repository, error) {
	start := time.Now()
	defer func() {
		needsWait := minRequestDuration - time.Since(start)
		if needsWait > 0 {
			log15.Debug("Waiting", "duration", needsWait)
			time.Sleep(needsWait)
		}
	}()
	repos, resp, err := gp.apiClient.Repositories.ListAll(&api.RepositoryListAllOptions{Since: since})
	if err != nil {
		return nil, err
	}
	gp.checkpoint = resp.NextPage
	gp.saveRepos(repos)
	if resp.Remaining < 100 {
		log15.Warn("Low remaining", "value", resp.Remaining)
	}

	return repos, nil
}

func (gp *githubProvider) getLastRepoId() (int, error) {
	result := api.Repository{}
	err := gp.repositoriesCollection.Find(nil).Sort("-_id").One(&result)
	if err == mgo.ErrNotFound {
		return 0, nil
	}

	return *result.ID, err
}

func (gp *githubProvider) saveRepos(repositories []*api.Repository) error {
	bulkOp := gp.repositoriesCollection.Bulk()
	for _, repo := range repositories {
		bulkOp.Insert(repo)
	}
	_, err := bulkOp.Run()

	return err
}
