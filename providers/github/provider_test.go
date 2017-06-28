package github

import (
	"database/sql"
	"errors"
	"io"
	"testing"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/github/model"

	. "gopkg.in/check.v1"
	ocore "gopkg.in/src-d/core.v0"
)

func Test(t *testing.T) {
	TestingT(t)
}

type GithubProviderSuite struct {
	DB       *sql.DB
	provider core.RepoProvider
}

var _ = Suite(&GithubProviderSuite{})

func (s *GithubProviderSuite) SetUpTest(c *C) {
	DB := ocore.Database()
	s.DB = DB

	err := core.DropTables(DB, providerName)
	c.Assert(err, IsNil)
	err = core.CreateGithubTable(DB)
	c.Assert(err, IsNil)

	s.provider = NewProvider(core.Config.Github.Token, s.DB)

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

	rs, err := model.NewRepositoryStore(s.DB).Find(model.NewRepositoryQuery())
	c.Assert(err, IsNil)
	repos, err := rs.All()
	c.Assert(err, IsNil)

	c.Assert(len(repos), Equals, 100)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart_ReposTwoPages(c *C) {
	for i := 0; i < 101; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}

	rs, err := model.NewRepositoryStore(s.DB).Find(model.NewRepositoryQuery())
	c.Assert(err, IsNil)
	repos, err := rs.All()
	c.Assert(err, IsNil)

	c.Assert(len(repos), Equals, 200)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_End(c *C) {
	repo := model.NewRepository()
	repo.GithubID = 99999999

	repos := []*model.Repository{
		repo,
	}

	// Simulate Ack
	githubProvider, ok := s.provider.(*provider)
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
