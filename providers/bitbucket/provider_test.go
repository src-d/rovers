package bitbucket

import (
	"errors"
	"io"
	"testing"

	"github.com/src-d/rovers/core"

	. "gopkg.in/check.v1"
)

const (
	testDatabase                       = "bitbucket-test"
	lastPage                           = "3000-01-00T17:25:17.038951+00:00"
	firstCheckpointWithGitRepositories = "2011-08-10T00:42:35.509559+00:00"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ProviderSuite struct {
	p core.RepoProvider
	c *core.Client
}

var _ = Suite(&ProviderSuite{c: core.NewClient(testDatabase)})

func (s *ProviderSuite) SetUpTest(c *C) {
	s.c.DropDatabase()
	s.p = NewProvider(testDatabase)
	bitbucketProvider, ok := s.p.(*provider)
	c.Assert(ok, Equals, true)
	bitbucketProvider.lastCheckpoint = firstCheckpointWithGitRepositories
}

func (s *ProviderSuite) TestProvider_Next(c *C) {
	r, err := s.p.Next()
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	result := bitbucketRepository{}
	err = s.c.Collection(repositoriesCollection).Find(nil).Sort("_id").One(&result)
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
