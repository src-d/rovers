package core

import "srcd.works/core.v0/model"

type RepoProvider interface {
	Next() (*model.Mention, error)
	Ack(error) error
	Close() error
	Name() string
}
