package commands

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
	"srcd.works/domain.v6/models/repository"
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
	repositoryRaw := &repository.Raw{
		Status:   repository.Initial,
		Provider: "test",
		URL:      "https://some.repo.url.com",
		IsFork:   true,
		VCS:      repository.Git,
	}

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

	obtainedRepositoryRaw := &repository.Raw{}
	err = job.Decode(obtainedRepositoryRaw)
	c.Assert(err, IsNil)
	testTime := time.Now()

	obtainedRepositoryRaw.CreatedAt = testTime
	obtainedRepositoryRaw.UpdatedAt = testTime
	repositoryRaw.CreatedAt = testTime
	repositoryRaw.UpdatedAt = testTime
	c.Assert(repositoryRaw, DeepEquals, obtainedRepositoryRaw)
}
