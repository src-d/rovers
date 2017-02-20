package commands

import (
	"github.com/src-d/rovers/core"

	ocore "srcd.works/core.v0"
)

type CmdCreateTables struct {
	CmdBase
}

func (c *CmdCreateTables) Execute(args []string) error {
	c.ChangeLogLevel()

	db := ocore.Database()

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
