package commands

import (
	"fmt"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	githubProviderName = "github"
)

var allowedProviders = []string{githubProviderName}

type CmdRrepoProviders struct {
	CmdBase
	Providers      []string `short:"p" long:"provider" optional:"no" description:"list of providers to execute."`
}

func (c *CmdRrepoProviders) Execute(args []string) error {
	c.ChangeLogLevel()

	for _, p := range c.Providers {
		switch p {
		case githubProviderName:
			log15.Info("Creating github provider")
		default:
			return fmt.Errorf("Provider '%s' not found. Allowed providers: %v",
				p, allowedProviders)
		}

	}

	return nil
}
