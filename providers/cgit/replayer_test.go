package cgit

import (
	"database/sql"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit/model"

	. "gopkg.in/check.v1"
	rcore "gopkg.in/src-d/core-retrieval.v0"
	kallax "gopkg.in/src-d/go-kallax.v1"
)

type CgitReplayerSuite struct {
	DB *sql.DB

	replayer core.RepoProvider
	store    *model.RepositoryStore
}

var _ = Suite(&CgitReplayerSuite{})

func (s *CgitReplayerSuite) SetUpTest(c *C) {
	db := rcore.Database()
	s.DB = db

	err := core.DropTables(db, core.CgitProviderName)
	c.Assert(err, IsNil)
	err = core.CreateCgitTables(db)
	c.Assert(err, IsNil)

	s.replayer = NewReplayer(db)
	s.store = model.NewRepositoryStore(db)

}

func (s *CgitReplayerSuite) TestNext(c *C) {
	s.createTuple(c, "foo", "bar", "baz")

	mention, err := s.replayer.Next()
	c.Assert(err, IsNil)
	c.Assert(mention.Endpoint, Equals, "foo")
	c.Assert(mention.IsFork, IsNil)
	c.Assert(len(mention.Aliases), Equals, 2)
	c.Assert(mention.Aliases, DeepEquals, []string{"bar", "baz"})

	_, err = s.replayer.Next()
	c.Assert(err, Equals, core.NoErrStopProvider)

	err = s.replayer.Close()
	c.Assert(err, IsNil)
}

func (s *CgitReplayerSuite) createTuple(c *C, url string, aliases ...string) {
	err := s.store.Insert(&model.Repository{ID: kallax.NewULID(), URL: url, Aliases: aliases})
	c.Assert(err, IsNil)
}
