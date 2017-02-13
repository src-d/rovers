package commands

import "github.com/src-d/rovers/core"

type CmdCreateTables struct {
	CmdBase
}

func (c *CmdCreateTables) Execute(args []string) error {
	c.ChangeLogLevel()

	db, err := core.NewDB()
	if err != nil {
		return err
	}

	err = core.CreateBitbucketTable(db)
	if err != nil {
		return err
	}

	err = core.CreateCgitTables(db)
	if err != nil {
		return err
	}

	return core.CreateGithubTable(db)
}
