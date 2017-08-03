package cgit

import (
	"database/sql"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit/model"

	rmodel "gopkg.in/src-d/core-retrieval.v0/model"
)

type replayer struct {
	store *model.RepositoryStore
	rs    *model.RepositoryResultSet
}

func NewReplayer(db *sql.DB) core.RepoProvider {
	return &replayer{store: model.NewRepositoryStore(db)}
}

func (r *replayer) Next() (*rmodel.Mention, error) {
	if r.rs == nil {
		if err := r.initialize(); err != nil {
			return nil, err
		}
	}

	if !r.rs.Next() {
		return nil, core.NoErrStopProvider
	}

	repository, err := r.rs.Get()
	if err != nil {
		return nil, err
	}

	return getMention(repository), nil
}

func (r *replayer) initialize() (err error) {
	r.rs, err = r.store.Find(model.NewRepositoryQuery())
	return
}

func (r *replayer) Ack(error) error {
	return nil
}

func (r *replayer) Close() error {
	return r.rs.Close()
}

func (r *replayer) Name() string {
	return core.CgitProviderName
}

func getMention(r *model.Repository) *rmodel.Mention {
	return &rmodel.Mention{
		Endpoint: r.URL,
		Provider: core.CgitProviderName,
		VCS:      rmodel.GIT,
		Aliases:  r.Aliases,
	}
}
