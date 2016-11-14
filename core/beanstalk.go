package core

import (
	"net"
	"time"

	"github.com/jpillora/backoff"
	"github.com/nutrun/lentil"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	maxDurationToRetry = 1 * time.Minute
	minDurationToRetry = 1 * time.Second

	writeOp = "write"
)

type beanstalkQueue struct {
	queue       *lentil.Beanstalkd
	connUrl     string
	name        string
	connBackoff *backoff.Backoff
	putBackoff  *backoff.Backoff
}

func NewBeanstalkQueue(connUrl string, name string) *beanstalkQueue {
	b := &beanstalkQueue{
		name:        name,
		connUrl:     connUrl,
		connBackoff: getBackoff(),
		putBackoff:  getBackoff(),
	}
	b.connect()

	return b
}

func getBackoff() *backoff.Backoff {
	return &backoff.Backoff{
		Jitter: true,
		Factor: 2,
		Max:    maxDurationToRetry,
		Min:    minDurationToRetry,
	}
}

func (b *beanstalkQueue) connect() {
	for {
		queue, err := lentil.Dial(b.connUrl)
		if err != nil {
			tts := b.connBackoff.Duration()
			log15.Error("beanstalk connection error",
				"error", err,
				"attempt", b.connBackoff.Attempt(),
				"time to sleep", tts)
			time.Sleep(tts)
			continue
		}
		b.connBackoff.Reset()
		queue.Use(b.name)
		b.queue = queue
		break
	}
}

func (b *beanstalkQueue) Put(
	body []byte,
	priority, delay, ttr int,
) uint64 {
	for {
		id, err := b.queue.Put(priority, delay, ttr, body)
		if err == nil {
			b.putBackoff.Reset()
			return id
		}
		switch specErr := err.(type) {
		case *net.OpError:
			if specErr.Op == writeOp {
				b.connect()
			}
		default:
			tts := b.putBackoff.Duration()
			log15.Error("beanstalk put error",
				"error", err,
				"attempt", b.putBackoff.Attempt(),
				"time to sleep", tts)
			time.Sleep(tts)
		}
	}
}

func (b *beanstalkQueue) QueueName() string {
	return b.name
}
