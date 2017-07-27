package github

import (
	"database/sql"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/github/model"

	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
)

type replayer struct {
	storer        *model.RepositoryStore
	rs            *model.RepositoryResultSet
	isInitialized bool
}

func NewReplayer(DB *sql.DB) core.RepoProvider {
	return &replayer{storer: model.NewRepositoryStore(DB)}
}

func (r *replayer) Next() (*rmodel.Mention, error) {
	if !r.isInitialized {
		if err := r.initialize(); err != nil {
			return nil, err
		}

		r.isInitialized = true
	}

	if !r.rs.Next() {
		return nil, core.NoErrStopProvider
	}

	repository, err := r.rs.Get()
	if err != nil {
		return nil, err
	}

	return getMention(repository.FullName, repository.Fork), nil
}

func (r *replayer) initialize() (err error) {
	r.rs, err = r.storer.Find(model.NewRepositoryQuery())
	return
}

func (r *replayer) Ack(error) error {
	return nil
}

func (r *replayer) Close() error {
	return r.rs.Close()
}

func (r *replayer) Name() string {
	return core.GithubProviderName
}
