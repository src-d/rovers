package commands

import "errors"

type CmdNoop struct {
}

func (cmd *CmdNoop) Execute(args []string) error {
	return errors.New("noop command called")
}
