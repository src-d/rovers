package commands

import (
	"github.com/src-d/rovers/core"

	"gopkg.in/src-d/core-retrieval.v0/model"
	"gopkg.in/src-d/framework.v0/queue"
)

type CmdQueue struct {
	Queue string `long:"queue" default:"rovers" description:"queue name"`
}

func (c *CmdQueue) getPersistFunction() (core.PersistFN, error) {
	broker, err := queue.NewBroker(core.Config.Broker.URL)
	if err != nil {
		return nil, err
	}

	q, err := broker.Queue(c.Queue)
	if err != nil {
		return nil, err
	}

	return func(repo *model.Mention) error {
		j, err := queue.NewJob()
		if err != nil {
			return err
		}

		if err = j.Encode(repo); err != nil {
			return err
		}

		return q.Publish(j)
	}, nil
}
