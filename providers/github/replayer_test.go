package github

import (
	"database/sql"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/github/model"

	. "gopkg.in/check.v1"
	rcore "gopkg.in/src-d/core-retrieval.v0"
)

type GithubReplayerSuite struct {
	DB *sql.DB

	replayer core.RepoProvider
	store    *model.RepositoryStore
}

var _ = Suite(&GithubReplayerSuite{})

func (s *GithubReplayerSuite) SetUpTest(c *C) {
	DB := rcore.Database()
	s.DB = DB

	err := core.DropTables(DB, core.GithubProviderName)
	c.Assert(err, IsNil)
	err = core.CreateGithubTable(DB)
	c.Assert(err, IsNil)

	s.replayer = NewReplayer(DB)
	s.store = model.NewRepositoryStore(DB)

}

func (s *GithubReplayerSuite) TestNext(c *C) {
	isFork := true
	s.createTuple(c, "some/repo", isFork)

	mention, err := s.replayer.Next()
	c.Assert(err, IsNil)
	c.Assert(mention.Endpoint, Equals, "git://github.com/some/repo")
	c.Assert(mention.IsFork, DeepEquals, &isFork)
	c.Assert(len(mention.Aliases), Equals, 3)
	c.Assert(mention.Aliases, DeepEquals, []string{
		"git://github.com/some/repo",
		"git@github.com:some/repo.git",
		"https://github.com/some/repo.git",
	})

	_, err = s.replayer.Next()
	c.Assert(err, Equals, core.NoErrStopProvider)

	err = s.replayer.Close()
	c.Assert(err, IsNil)
}

func (s *GithubReplayerSuite) createTuple(c *C, name string, fork bool) {
	err := s.store.Insert(&model.Repository{FullName: name, Fork: fork})
	c.Assert(err, IsNil)
}
