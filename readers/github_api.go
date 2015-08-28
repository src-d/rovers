package readers

import (
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/inconshreveable/log15.v2"

	api "github.com/mcuadros/go-github/github"
)

const (
	GithubToken        = "b286be1a91d5656483209a9f3fdf120ab1174b67"
	MinRequestDuration = time.Hour / 5000
)

type GithubAPI struct {
	client *api.Client
}

func NewGithubAPI() *GithubAPI {
	token := &oauth2.Token{AccessToken: GithubToken}
	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	return &GithubAPI{api.NewClient(client)}
}

func (g *GithubAPI) GetAllRepositories(since int) ([]api.Repository, *api.Response, error) {
	start := time.Now()
	defer func() {
		needsWait := MinRequestDuration - time.Since(start)
		if needsWait > 0 {
			log15.Info("Waiting", "duration", needsWait)
			time.Sleep(needsWait)
		}
	}()

	options := &api.RepositoryListAllOptions{Since: since}
	repos, resp, err := g.client.Repositories.ListAll(options)
	if err != nil {
		return nil, resp, err
	}

	if resp.Remaining < 100 {
		log15.Info("Low remaining", "value", resp.Remaining)
	}

	return repos, resp, nil
}

func (g *GithubAPI) GetAllUsers(since int) ([]api.User, *api.Response, error) {
	start := time.Now()
	defer func() {
		needsWait := MinRequestDuration - time.Since(start)
		if needsWait > 0 {
			log15.Info("Waiting", "duration", needsWait)
			time.Sleep(needsWait)
		}
	}()

	o := &api.UserListOptions{Since: since}
	users, resp, err := g.client.Users.ListAll(o)
	if err != nil {
		return nil, resp, err
	}

	if resp.Remaining < 100 {
		log15.Info("Low remaining", "value", resp.Remaining)
	}

	return users, resp, nil
}
