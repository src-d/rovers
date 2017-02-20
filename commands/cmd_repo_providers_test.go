package commands

import (
	"testing"

	"github.com/src-d/rovers/core"

	. "gopkg.in/check.v1"
	"srcd.works/core.v0/model"
	"srcd.works/framework.v0/queue"
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
		Queue: "test",
	}
}

func (s *CmdRepoProviderSuite) TestCmdRepoProvider_getPersistFunction_CorrectlySerialized(c *C) {
	repositoryRaw := &model.Mention{
		Provider: "test",
		Endpoint: "https://some.repo.url.com",
		VCS:      model.GIT,
		Context:  make(map[string]string),
	}

	repositoryRaw.Context["test"] = "bla"

	f, err := s.cmdProviders.getPersistFunction()
	c.Assert(err, IsNil)
	err = f(repositoryRaw)
	c.Assert(err, IsNil)

	broker, err := queue.NewBroker(core.Config.Broker.URL)
	c.Assert(err, IsNil)
	queue, err := broker.Queue(s.cmdProviders.Queue)
	c.Assert(err, IsNil)
	jobIter, err := queue.Consume()
	c.Assert(err, IsNil)

	job, err := jobIter.Next()
	c.Assert(err, IsNil)

	obtainedRepositoryRaw := &model.Mention{}
	err = job.Decode(obtainedRepositoryRaw)
	c.Assert(err, IsNil)

	// TODO Duration types are not serialized correctly
	obtainedRepositoryRaw.CreatedAt = repositoryRaw.CreatedAt
	obtainedRepositoryRaw.UpdatedAt = repositoryRaw.UpdatedAt
	c.Assert(repositoryRaw, DeepEquals, obtainedRepositoryRaw)
}
