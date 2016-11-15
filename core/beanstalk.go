package core

import (
	"time"

	"github.com/ajnavarro/beanstalk"
)

const (
	tcpNetwork = "tcp"
)

type beanstalkQueue struct {
	conn *beanstalk.Conn
	name string

	tube    *beanstalk.Tube
	tubeSet *beanstalk.TubeSet
}

func NewBeanstalkQueue(addr string, retries int, delay time.Duration, name string) (*beanstalkQueue, error) {
	conn, err := beanstalk.Dial(&beanstalk.Config{
		Addr:    addr,
		Delay:   delay,
		Network: tcpNetwork,
		Retries: retries,
	})
	if err != nil {
		return nil, err
	}

	return &beanstalkQueue{
		conn: conn,
		name: name,

		tube:    &beanstalk.Tube{conn, name},      // Put
		tubeSet: beanstalk.NewTubeSet(conn, name), // Reserve
	}, nil
}

func (b *beanstalkQueue) Put(
	body []byte,
	priority uint32,
	delay, ttr time.Duration,
) (id uint64, err error) {

	return b.tube.Put(body, priority, delay, ttr)
}

func (b *beanstalkQueue) QueueName() string {
	return b.name
}

func (b *beanstalkQueue) Bury(id uint64, priority uint32) error {
	return b.conn.Bury(id, priority)
}

func (b *beanstalkQueue) Delete(id uint64) error {
	return b.conn.Delete(id)
}

func (b *beanstalkQueue) Reserve(timeout time.Duration) (id uint64, body []byte, err error) {
	return b.tubeSet.Reserve(timeout)
}
