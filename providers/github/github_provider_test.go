package github

import (
	"errors"
	"io"
	"testing"

	"github.com/src-d/rovers/core"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type GithubProviderSuite struct {
	client *core.Client
}

var _ = Suite(&GithubProviderSuite{})

func (s *GithubProviderSuite) SetUpTest(c *C) {
	s.client = core.NewClient()
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart(c *C) {
	s.client.DropDatabase()
	provider := NewProvider(&GithubConfig{})
	for i := 0; i < 101; i++ {
		repoUrl, err := provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, Not(Equals), "")
		err = provider.Ack(nil)
		c.Assert(err, IsNil)
	}
}

func (s *GithubProviderSuite) TestGithubProvider_Next_End(c *C) {
	s.client.DropDatabase()
	provider := NewProvider(&GithubConfig{})
	provider.setCheckpoint(&GithubData{
		Checkpoint: 99999999,
	})
	repoUrl, err := provider.Next()
	c.Assert(repoUrl, Equals, "")
	c.Assert(err, Equals, io.EOF)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_Retry(c *C) {
	s.client.DropDatabase()
	provider := NewProvider(&GithubConfig{})
	repoUrl, err := provider.Next()
	c.Assert(err, IsNil)
	c.Assert(repoUrl, Not(Equals), "")

	// Simulate an error
	provider.Ack(errors.New("WOOPS"))
	repoUrl2, err := provider.Next()

	c.Assert(err, IsNil)
	c.Assert(repoUrl, Not(Equals), "")
	c.Assert(repoUrl, Equals, repoUrl2)
}
