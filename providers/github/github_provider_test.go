package github

import (
	"errors"
	"io"
	"testing"

	"github.com/src-d/go-github/github"
	"github.com/src-d/rovers/core"
	. "gopkg.in/check.v1"
)

const testDatabase = "github-test"

func Test(t *testing.T) {
	TestingT(t)
}

type GithubProviderSuite struct {
	client   *core.Client
	provider core.RepoProvider
}

var _ = Suite(&GithubProviderSuite{
	client: core.NewClient(testDatabase),
})

func (s *GithubProviderSuite) SetUpTest(c *C) {
	s.client.DropDatabase()
	config := &GithubConfig{GithubToken: core.Config.Github.Token, Database: testDatabase}
	s.provider = NewProvider(config)

}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart(c *C) {
	for i := 0; i < 101; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart_Repos(c *C) {
	for i := 0; i < 100; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}

	res := []github.Repository{}
	err := s.client.Collection(repositoryCollection).Find(nil).All(&res)
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 100)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart_ReposTwoPages(c *C) {
	for i := 0; i < 101; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}

	res := []github.Repository{}
	err := s.client.Collection(repositoryCollection).Find(nil).All(&res)
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 200)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_End(c *C) {
	lastRepoId := 99999999
	repos := []*github.Repository{{
		ID: &lastRepoId,
	}}

	// Simulate Ack
	githubProvider, ok := s.provider.(*githubProvider)
	c.Assert(ok, Equals, true)
	err := githubProvider.saveRepos(repos)
	c.Assert(err, IsNil)

	repoUrl, err := s.provider.Next()
	c.Assert(repoUrl, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_Retry(c *C) {
	repoUrl, err := s.provider.Next()
	c.Assert(err, IsNil)
	c.Assert(repoUrl, NotNil)

	// Simulate an error
	s.provider.Ack(errors.New("WOOPS"))
	repoUrl2, err := s.provider.Next()

	c.Assert(err, IsNil)
	c.Assert(repoUrl, NotNil)
	c.Assert(repoUrl, DeepEquals, repoUrl2)
}
