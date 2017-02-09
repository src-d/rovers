package core

import "srcd.works/core.v0/models"

type RepoProvider interface {
	Next() (*models.Mention, error)
	Ack(error) error
	Close() error
	Name() string
}
