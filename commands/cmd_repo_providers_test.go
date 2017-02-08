package commands

import (
	"testing"

	. "gopkg.in/check.v1"
	"srcd.works/core.v0/models"
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
		Queue:  "test",
		Broker: "amqp://guest:guest@localhost:5672/",
	}
}

func (s *CmdRepoProviderSuite) TestCmdRepoProvider_getPersistFunction_CorrectlySerialized(c *C) {
	repositoryRaw := &models.Mention{
		Provider: "test",
		Endpoint: "https://some.repo.url.com",
		VCS:      models.GIT,
		Context:  make(map[string]string),
	}

	repositoryRaw.Context["test"] = "bla"

	f, err := s.cmdProviders.getPersistFunction()
	c.Assert(err, IsNil)
	err = f(repositoryRaw)
	c.Assert(err, IsNil)

	broker, err := queue.NewBroker(s.cmdProviders.Broker)
	c.Assert(err, IsNil)
	queue, err := broker.Queue(s.cmdProviders.Queue)
	c.Assert(err, IsNil)
	jobIter, err := queue.Consume()
	c.Assert(err, IsNil)

	job, err := jobIter.Next()
	c.Assert(err, IsNil)

	obtainedRepositoryRaw := &models.Mention{}
	err = job.Decode(obtainedRepositoryRaw)
	c.Assert(err, IsNil)

	// TODO Duration types are not serialized correctly
	obtainedRepositoryRaw.CreatedAt = repositoryRaw.CreatedAt
	obtainedRepositoryRaw.UpdatedAt = repositoryRaw.UpdatedAt
	c.Assert(repositoryRaw, DeepEquals, obtainedRepositoryRaw)
}
