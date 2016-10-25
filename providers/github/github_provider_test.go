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
	client   *core.Client
	provider *githubProvider
}

var _ = Suite(&GithubProviderSuite{
	client: core.NewClient(providerName),
})

func (s *GithubProviderSuite) SetUpTest(c *C) {
	s.client.DropDatabase()
	s.provider = NewProvider(&GithubConfig{})
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

	res := []githubData{}
	err := s.client.Collection(providerName).Find(nil).All(&res)
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert(len(res[0].Repositories), Equals, 100)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart_ReposTwoPages(c *C) {
	for i := 0; i < 101; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}

	res := []githubData{}
	err := s.client.Collection(providerName).Find(nil).All(&res)
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert(len(res[0].Repositories), Equals, 100)
	c.Assert(len(res[1].Repositories), Equals, 100)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_End(c *C) {
	s.provider.setCheckpoint(&githubData{
		Checkpoint: 99999999,
	})
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
