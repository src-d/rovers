package bitbucket

import (
	"github.com/pkg/errors"
	"github.com/src-d/rovers/core"
	. "gopkg.in/check.v1"
	"io"
)

const testDatabase = "bitbucket-test"
const lastPage = "3000-01-00T17:25:17.038951+00:00"

type ProviderSuite struct {
	p core.RepoProvider
	c *core.Client
}

var _ = Suite(&ProviderSuite{c: core.NewClient(providerName)})

func (s *ProviderSuite) SetUpTest(c *C) {
	s.c.DropDatabase()
	s.p = NewProvider(testDatabase)
}

func (s *ProviderSuite) TestProvider_Next(c *C) {
	r, err := s.p.Next()
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
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
