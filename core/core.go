package core

type RepoProvider interface {
	Next() (string, error)
	Ack(error) error
	Close() error
	Name() string
}
