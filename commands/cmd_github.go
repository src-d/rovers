package commands

import (
	"gopkg.in/inconshreveable/log15.v2"
)

type CmdGitHub struct {
	CmdBase
}

func (c *CmdGitHub) Execute(args []string) error {
	log15.Warn("This subcommand does nothing")
	return nil
}
