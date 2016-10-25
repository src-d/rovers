package core

import "gop.kg/src-d/domain@v6/models/repository"

type RepoProvider interface {
	Next() (*repository.Raw, error)
	Ack(error) error
	Close() error
	Name() string
}
