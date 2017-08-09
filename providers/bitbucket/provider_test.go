package bitbucket

import (
	"database/sql"
	"errors"
	"io"
	"testing"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/bitbucket/model"

	. "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
	rcore "gopkg.in/src-d/core-retrieval.v0"
	"gopkg.in/src-d/go-kallax.v1"
)

const (
	lastPage                           = "3000-01-00T17:25:17.038951+00:00"
	firstCheckpointWithGitRepositories = "2011-08-10T00:42:35.509559+00:00"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ProviderSuite struct {
	p  core.RepoProvider
	DB *sql.DB
}

var _ = Suite(&ProviderSuite{})

func (s *ProviderSuite) SetUpTest(c *C) {
	httpmock.Activate()
	LoadAssets(c)

	DB := rcore.Database()
	s.DB = DB

	err := core.DropTables(DB, core.BitbucketProviderName)
	c.Assert(err, IsNil)

	err = core.CreateBitbucketTable(DB)
	c.Assert(err, IsNil)

	s.p = NewProvider(s.DB)
	bitbucketProvider, ok := s.p.(*provider)
	c.Assert(ok, Equals, true)
	bitbucketProvider.lastCheckpoint = firstCheckpointWithGitRepositories
}

func (s *ProviderSuite) TearDownTest(c *C) {
	httpmock.DeactivateAndReset()
}

func (s *ProviderSuite) TestProvider_Next(c *C) {
	r, err := s.p.Next()
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	result, err := model.NewRepositoryStore(s.DB).FindOne(
		model.NewRepositoryQuery().
			Order(kallax.Asc(model.Schema.Repository.CreatedAt)),
	)

	c.Assert(err, IsNil)
	c.Assert(result.Links.Clone[0].Href, Equals, r.Endpoint)
}

func (s *ProviderSuite) TestProvider_NextLast(c *C) {
	bitbucketProvider, ok := s.p.(*provider)
	c.Assert(ok, Equals, true)
	bitbucketProvider.lastCheckpoint = lastPage
	_, err := s.p.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ProviderSuite) TestProvider_NextRetry(c *C) {
	r, err := s.p.Next()
	c.Assert(r, NotNil)
	c.Assert(err, IsNil)
	err = s.p.Ack(errors.New("WOOPS"))
	c.Assert(err, IsNil)
	r2, err := s.p.Next()
	c.Assert(err, IsNil)
	c.Assert(r, DeepEquals, r2)
}

func (s *ProviderSuite) TestProvider_NextNoAck(c *C) {
	r, err := s.p.Next()
	c.Assert(r, NotNil)
	c.Assert(err, IsNil)
	r2, err := s.p.Next()
	c.Assert(r2, NotNil)
	c.Assert(err, IsNil)
	c.Assert(r, DeepEquals, r2)
}

func (s *ProviderSuite) TestProvider_NextAckNext(c *C) {
	r, err := s.p.Next()
	c.Assert(r, NotNil)
	c.Assert(err, IsNil)
	err = s.p.Ack(nil)
	c.Assert(err, IsNil)
	r2, err := s.p.Next()
	c.Assert(r2, NotNil)
	c.Assert(err, IsNil)
	c.Assert(r, Not(DeepEquals), r2)
}
