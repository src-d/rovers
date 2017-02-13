package core

import (
	"errors"
	"io"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
	"srcd.works/core.v0/model"
)

const (
	maxRetries            = 3
	secondsBetweenRetries = 10
)

type PersistFN func(*model.Mention) error

var errBadAck = errors.New("error while executing ack")

type Watcher struct {
	providers      []RepoProvider
	persist        PersistFN
	timeToSleep    time.Duration
	timeToRetryAck time.Duration
}

func NewWatcher(providers []RepoProvider, persist PersistFN,
	timeToSleep time.Duration, timeToRetryAck time.Duration) *Watcher {
	if timeToRetryAck == 0 {
		timeToRetryAck = time.Second * secondsBetweenRetries
	}
	return &Watcher{
		providers:      providers,
		persist:        persist,
		timeToSleep:    timeToSleep,
		timeToRetryAck: timeToRetryAck,
	}
}

func (w *Watcher) Start() {
	for _, provider := range w.providers {
		go func(p RepoProvider) {
			for {
				if err := w.handleProviderResult(p); err == errBadAck {
					break
				}
			}
		}(provider)
	}
}

func (w *Watcher) handleProviderResult(p RepoProvider) error {
	mention, err := p.Next()
	switch err {
	case io.EOF:
		log15.Info("no more repositories, "+
			"waiting for more...",
			"time to sleep", w.timeToSleep)
		time.Sleep(w.timeToSleep)
	case nil:
		log15.Info("getting new repository", "provider", p.Name(), "repository", mention.Endpoint)
		err := w.persist(mention)
		if err != nil {
			log15.Error("error saving new repo", "error", err, "repository", mention.Endpoint)
		}
		retries := 0
		for retries != maxRetries {
			ackErr := p.Ack(err)
			if ackErr != nil {
				log15.Error("error setting ack", "ack error", ackErr)
				retries++
				time.Sleep(w.timeToRetryAck)
			} else {
				break
			}
		}
		if retries == maxRetries {
			log15.Error("error in ack. Shutting down provider",
				"provider", p.Name(), "retries", retries, "time to retry", w.timeToRetryAck)
			p.Close()
			return errBadAck
		}
	default:
		log15.Error("error obtaining new repository", "provider", p.Name(), "error", err)
	}

	return nil
}
