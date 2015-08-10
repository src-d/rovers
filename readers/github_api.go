package readers

import (
	"fmt"
	"time"

	"code.google.com/p/goauth2/oauth"
	api "github.com/mcuadros/go-github/github"
	"github.com/tyba/oss/sources/social/http"
)

var MinRequestDuration = time.Hour / 5000

type GithubAPIReader struct {
	client *api.Client
}

func NewGithubAPIReader(client *http.Client) *GithubAPIReader {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: "e568ba2365b2bc198da8c5571a4cfb99830bb5ed"},
	}

	return &GithubAPIReader{api.NewClient(t.Client())}
}

func (g *GithubAPIReader) GetAllRepositories(since int) ([]api.Repository, *api.Response, error) {
	start := time.Now()
	defer func() {
		needsWait := MinRequestDuration - time.Now().Sub(start)
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

func (g *GithubAPIReader) GetAllUsers(since int) ([]api.User, *api.Response, error) {
	start := time.Now()
	defer func() {
		needsWait := MinRequestDuration - time.Now().Sub(start)
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
