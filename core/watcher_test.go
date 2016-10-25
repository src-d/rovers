package core

import (
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"gop.kg/src-d/domain@v6/models/repository"
	. "gopkg.in/check.v1"
	"gopkg.in/inconshreveable/log15.v2"
)

var persistFunction = func(rawRepo *repository.Raw) error {
	log15.Debug("Persisting new url", "repoUrl", rawRepo)
	return nil
}

func Test(t *testing.T) {
	TestingT(t)
}

type WatcherSuite struct {
	rawRepo *repository.Raw
}

var _ = Suite(&WatcherSuite{
	rawRepo: &repository.Raw{
		IsFork: true,
		URL:    "SOME_REPO",
		Status: repository.Initial,
	},
})

func (s *WatcherSuite) TestWatcher_EOF(c *C) {
	provider := &EOFProvider{Repo: s.rawRepo}
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
		Repo:      s.rawRepo,
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
	Repo                  *repository.Raw
}

func (p *EOFProvider) Next() (*repository.Raw, error) {
	mutex.Lock()
	defer mutex.Unlock()

	p.NumberOfCalls++
	switch p.NumberOfCalls {
	case 1:
		return p.Repo, nil
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
