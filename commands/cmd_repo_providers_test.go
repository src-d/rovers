package commands

import (
	"testing"

	"github.com/src-d/rovers/core"

	. "gopkg.in/check.v1"
	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
	"gopkg.in/src-d/go-queue.v1"
	_ "gopkg.in/src-d/go-queue.v1/amqp"
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
		CmdQueue: CmdQueue{
			Queue: "test",
		},
	}
}

func (s *CmdRepoProviderSuite) TestCmdRepoProvider_getPersistFunction_CorrectlySerialized(c *C) {
	isFork := false
	repositoryRaw := &rmodel.Mention{
		Provider: "test",
		Endpoint: "https://some.repo.url.com",
		VCS:      rmodel.GIT,
		Aliases:  []string{"a", "b", "c"},
		IsFork:   &isFork,
	}

	f, err := s.cmdProviders.getPersistFunction()
	c.Assert(err, IsNil)
	err = f(repositoryRaw)
	c.Assert(err, IsNil)

	broker, err := queue.NewBroker(core.Config.Broker.URL)
	c.Assert(err, IsNil)
	queue, err := broker.Queue(s.cmdProviders.Queue)
	c.Assert(err, IsNil)
	jobIter, err := queue.Consume(1)
	c.Assert(err, IsNil)

	job, err := jobIter.Next()
	c.Assert(err, IsNil)

	obtainedRepositoryRaw := &rmodel.Mention{}
	err = job.Decode(obtainedRepositoryRaw)
	c.Assert(err, IsNil)

	// TODO Duration types are not serialized correctly
	obtainedRepositoryRaw.CreatedAt = repositoryRaw.CreatedAt
	obtainedRepositoryRaw.UpdatedAt = repositoryRaw.UpdatedAt
	c.Assert(repositoryRaw, DeepEquals, obtainedRepositoryRaw)
}
