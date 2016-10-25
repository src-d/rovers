package core

import (
	"errors"
	"io"
	"time"

	"gop.kg/src-d/domain@v6/models/repository"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	maxRetries            = 3
	secondsBetweenRetries = 10
)

var errBadAck = errors.New("Error while executing ACK")

type Watcher struct {
	providers      []RepoProvider
	persist        func(*repository.Raw) error
	timeToSleep    time.Duration
	timeToRetryAck time.Duration
}

func NewWatcher(providers []RepoProvider, persist func(*repository.Raw) error,
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
	repoRaw, err := p.Next()
	switch err {
	case io.EOF:
		log15.Info("No more repositories, "+
			"waiting for more...",
			"time to sleep", w.timeToSleep)
		time.Sleep(w.timeToSleep)
	case nil:
		log15.Info("Getting new repository", "provider", p.Name(), "repository", repoRaw.URL)
		err := w.persist(repoRaw)
		if err != nil {
			log15.Error("Error saving new repo", "error", err, "repository", repoRaw.URL)
		}
		retries := 0
		for retries != maxRetries {
			ackErr := p.Ack(err)
			if ackErr != nil {
				log15.Error("Error setting Ack", "ackErr", ackErr)
				retries++
				time.Sleep(w.timeToRetryAck)
			} else {
				break
			}
		}
		if retries == maxRetries {
			log15.Error("Error in ACK. Shutting down provider.",
				"provider", p.Name(), "retries", retries, "timeToRetry", w.timeToRetryAck)
			p.Close()
			return errBadAck
		}
	default:
		log15.Error("Error obtaining new repo...", "error", err)
	}

	return nil
}
