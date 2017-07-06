package commands

import (
	"github.com/src-d/rovers/core"

	rcore "gopkg.in/src-d/core-retrieval.v0"
)

type CmdCreateTables struct {
	CmdBase
}

func (c *CmdCreateTables) Execute(args []string) error {
	c.ChangeLogLevel()

	db := rcore.Database()

	err := core.CreateBitbucketTable(db)
	if err != nil {
		return err
	}

	err = core.CreateCgitTables(db)
	if err != nil {
		return err
	}

	return core.CreateGithubTable(db)
}
