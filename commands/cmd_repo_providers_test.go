package commands

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/kr/beanstalk"
	"github.com/sourcegraph/go-vcsurl"
	"github.com/src-d/rovers/core"
	"gop.kg/src-d/domain@v6/models/repository"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type CmdRepoProviderSuite struct {
	cmdProviders *CmdRepoProviders
}

var _ = Suite(&CmdRepoProviderSuite{})

func (s *CmdRepoProviderSuite) SetUpTest(c *C) {
	s.cmdProviders = &CmdRepoProviders{
		QueueName: "test",
		Beanstalk: "127.0.0.1:11300",
	}
}

func (s *CmdRepoProviderSuite) TestCmdRepoProvider_getPersistFunction_CorrectlySerialized(c *C) {
	repositoryRaw := &repository.Raw{
		Status:   repository.Initial,
		Provider: "test",
		URL:      "https://some.repo.url.com",
		IsFork:   true,
		VCS:      vcsurl.Git,
	}

	f, err := s.cmdProviders.getPersistFunction()
	c.Assert(err, IsNil)
	err = f(repositoryRaw)
	c.Assert(err, IsNil)

	conn, err := beanstalk.Dial("tcp", s.cmdProviders.Beanstalk)
	c.Assert(err, IsNil)
	queue := core.NewBeanstalkQueue(conn, s.cmdProviders.QueueName)
	_, body, err := queue.Reserve(500 * time.Millisecond)
	c.Assert(err, IsNil)

	obtainedRepositoryRaw := &repository.Raw{}
	err = gob.NewDecoder(bytes.NewReader(body)).Decode(obtainedRepositoryRaw)
	c.Assert(err, IsNil)
	c.Assert(repositoryRaw, DeepEquals, obtainedRepositoryRaw)
}
