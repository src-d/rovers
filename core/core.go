package core

import "srcd.works/domain.v6/models/repository"

type RepoProvider interface {
	Next() (*repository.Raw, error)
	Ack(error) error
	Close() error
	Name() string
}
