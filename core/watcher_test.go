package core

import (
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	. "gopkg.in/check.v1"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/src-d/core-retrieval.v0/model"
)

var persistFunction = func(rawRepo *model.Mention) error {
	log15.Debug("Persisting new url", "repoUrl", rawRepo)
	return nil
}

func Test(t *testing.T) {
	TestingT(t)
}

type WatcherSuite struct {
	mention *model.Mention
}

var _ = Suite(&WatcherSuite{
	mention: &model.Mention{
		Endpoint: "SOME_REPO",
		Provider: "SOME_PROVIDER",
	},
})

func (s *WatcherSuite) TestWatcher_EOF(c *C) {
	provider := &EOFProvider{Mention: s.mention}
	providers := []RepoProvider{provider}
	watcher := NewWatcher(providers, persistFunction, time.Second, time.Second)
	watcher.Start()
	time.Sleep(time.Second * 5)
	gt := provider.NumberOfCalls >= 5
	c.Assert(gt, Equals, true)
}

func (s *WatcherSuite) TestWatcher_AckRetries(c *C) {
	provider := &EOFProvider{
		FailOnAck: true,
		Mention:   s.mention,
	}
	providers := []RepoProvider{provider}
	watcher := NewWatcher(providers, persistFunction, time.Second, time.Second)
	watcher.Start()
	time.Sleep(time.Second * 5)
	c.Assert(provider.NumberOfCalls, Equals, 1)
	c.Assert(provider.NumberOfAckErrorCalls, Equals, 3)
}

var (
	mutex *sync.Mutex = &sync.Mutex{}
)

type EOFProvider struct {
	NumberOfCalls         int
	NumberOfAckErrorCalls int
	FailOnAck             bool
	Mention               *model.Mention
}

func (p *EOFProvider) Next() (*model.Mention, error) {
	mutex.Lock()
	defer mutex.Unlock()

	p.NumberOfCalls++
	switch p.NumberOfCalls {
	case 1:
		return p.Mention, nil
	case 2:
		return nil, errors.New("OTHER ERROR")
	default:
		return nil, io.EOF
	}
}

func (p *EOFProvider) Ack(err error) error {
	mutex.Lock()
	defer mutex.Unlock()
	if p.FailOnAck {
		p.NumberOfAckErrorCalls++
		return errors.New("SOME ACK ERROR")
	}
	return nil
}

func (p *EOFProvider) Close() error {
	return nil
}

func (p *EOFProvider) Name() string {
	return "EOFProvider"
}
