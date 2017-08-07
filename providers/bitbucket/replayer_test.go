package bitbucket

import (
	"database/sql"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/bitbucket/model"

	. "gopkg.in/check.v1"
	rcore "gopkg.in/src-d/core-retrieval.v0"
	kallax "gopkg.in/src-d/go-kallax.v1"
)

type BitbucketReplayerSuite struct {
	DB *sql.DB

	replayer core.RepoProvider
	store    *model.RepositoryStore
}

var _ = Suite(new(BitbucketReplayerSuite))

func (s *BitbucketReplayerSuite) SetUpTest(c *C) {
	db := rcore.Database()
	s.DB = db

	err := core.DropTables(s.DB, core.BitbucketProviderName)
	c.Assert(err, IsNil)

	err = core.CreateBitbucketTable(db)
	c.Assert(err, IsNil)

	s.replayer = NewReplayer(db)
	s.store = model.NewRepositoryStore(db)
}

func (s *BitbucketReplayerSuite) TestNext(c *C) {
	s.createTuple(
		c,
		true,
		link{"https://foo.bar", httpsCloneKey},
		link{"ssh://foo.bar", "ssh"},
	)
	s.createTuple(
		c,
		false,
		link{"git://bar.baz", "git"},
		link{"ssh://bar.baz", "ssh"},
	)

	mention, err := s.replayer.Next()
	c.Assert(err, IsNil)
	c.Assert(mention.Endpoint, Equals, "https://foo.bar")
	c.Assert(*mention.IsFork, Equals, true)
	c.Assert(len(mention.Aliases), Equals, 2)
	c.Assert(mention.Aliases, DeepEquals, []string{
		"https://foo.bar",
		"ssh://foo.bar",
	})

	mention, err = s.replayer.Next()
	c.Assert(err, IsNil)
	c.Assert(mention.Endpoint, Equals, "git://bar.baz")
	c.Assert(*mention.IsFork, Equals, false)
	c.Assert(len(mention.Aliases), Equals, 2)
	c.Assert(mention.Aliases, DeepEquals, []string{
		"git://bar.baz",
		"ssh://bar.baz",
	})

	_, err = s.replayer.Next()
	c.Assert(err, Equals, core.NoErrStopProvider)

	err = s.replayer.Close()
	c.Assert(err, IsNil)
}

type link struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

func (s *BitbucketReplayerSuite) createTuple(c *C, fork bool, links ...link) {
	var parent *model.Parent
	if fork {
		parent = &model.Parent{UUID: "foo"}
	}

	repo := &model.Repository{
		ID:     kallax.NewULID(),
		Parent: parent,
	}

	for _, l := range links {
		repo.Links.Clone = append(repo.Links.Clone, l)
	}

	err := s.store.Insert(repo)
	c.Assert(err, IsNil)
}
