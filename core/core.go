package core

import "srcd.works/core-retrieval.v0/model"

type RepoProvider interface {
	Next() (*model.Mention, error)
	Ack(error) error
	Close() error
	Name() string
}
