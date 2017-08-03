package commands

import (
	"fmt"
	"time"

	"github.com/src-d/rovers/core"

	"github.com/jessevdk/go-flags"
	"github.com/src-d/rovers/providers/cgit"
	"github.com/src-d/rovers/providers/github"
	rcore "gopkg.in/src-d/core-retrieval.v0"
)

type CmdReplay struct {
	CmdBase
	CmdQueue

	Providers []string `short:"p" long:"provider" optional:"yes" description:"list of providers to execute."`
}

func (c *CmdReplay) Execute(args []string) error {
	c.ChangeLogLevel()

	if len(c.Providers) == 0 {
		return &flags.Error{Type: flags.ErrHelp, Message: "ERROR: no replayers to execute provided"}
	}

	DB := rcore.Database()

	var replayers []core.RepoProvider
	for _, rep := range c.Providers {
		switch rep {
		case core.GithubProviderName:
			replayers = append(replayers, github.NewReplayer(DB))
		case core.CgitProviderName:
			replayers = append(replayers, cgit.NewReplayer(DB))
		default:
			return &flags.Error{
				Type:    flags.ErrInvalidChoice,
				Message: fmt.Sprintf("ERROR: provider %s not supported for replay", rep),
			}
		}
	}

	f, err := c.getPersistFunction()
	if err != nil {
		return err
	}
	watcher := core.NewWatcher(replayers, f, 24*time.Hour, time.Second*15)
	watcher.Start()

	select {}

	return nil
}
