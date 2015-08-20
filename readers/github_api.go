package readers

import (
	"fmt"
	"time"

	"code.google.com/p/goauth2/oauth"
	api "github.com/mcuadros/go-github/github"
)

var MinRequestDuration = time.Hour / 5000

type GithubAPI struct {
	client *api.Client
}

func NewGithubAPI() *GithubAPI {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: "e568ba2365b2bc198da8c5571a4cfb99830bb5ed"},
	}

	return &GithubAPI{api.NewClient(t.Client())}
}

func (g *GithubAPI) GetAllRepositories(since int) ([]api.Repository, *api.Response, error) {
	start := time.Now()
	defer func() {
		needsWait := MinRequestDuration - time.Since(start)
		if needsWait > 0 {
			fmt.Println("waiting ", needsWait)
			time.Sleep(needsWait)
		}
	}()

	o := &api.RepositoryListAllOptions{Since: since}
	repos, resp, err := g.client.Repositories.ListAll(o)
	if err != nil {
		return nil, resp, err
	}

	if resp.Remaining < 100 {
		fmt.Println("low remaining", resp.Remaining)
	}

	return repos, resp, nil
}

func (g *GithubAPI) GetAllUsers(since int) ([]api.User, *api.Response, error) {
	start := time.Now()
	defer func() {
		needsWait := MinRequestDuration - time.Since(start)
		if needsWait > 0 {
			fmt.Println("waiting ", needsWait)
			time.Sleep(needsWait)
		}
	}()

	o := &api.UserListOptions{Since: since}
	users, resp, err := g.client.Users.ListAll(o)
	if err != nil {
		return nil, resp, err
	}

	if resp.Remaining < 100 {
		fmt.Println("low remaining", resp.Remaining)
	}

	return users, resp, nil
}
