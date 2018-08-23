package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/bitbucket"
	"github.com/src-d/rovers/providers/cgit"
	"github.com/src-d/rovers/providers/github"

	"gopkg.in/inconshreveable/log15.v2"
	rcore "gopkg.in/src-d/core-retrieval.v0"
)

var allowedProviders = []string{core.GithubProviderName, core.CgitProviderName, core.BitbucketProviderName}

type CmdRepoProviders struct {
	CmdBase
	CmdQueue

	Providers   []string      `short:"p" long:"provider" optional:"yes" description:"provider to execute, any of [github bitbucket cgit]. (If you don't set any provider, all supported provider will be used)"`
	WatcherTime time.Duration `short:"t" long:"watcher-time" optional:"no" default:"1h" description:"Time to try again to get new repos"`
}

func (c *CmdRepoProviders) Execute(args []string) error {
	c.ChangeLogLevel()

	if len(c.Providers) == 0 {
		log15.Info("No providers added using --provider option. Executing all known providers",
			"providers", allowedProviders)
		c.Providers = allowedProviders
	}

	DB := rcore.Database()

	providers := []core.RepoProvider{}
	for _, p := range c.Providers {
		switch p {
		case core.GithubProviderName:
			log15.Info("Creating github provider")
			if core.Config.Github.Token == "" {
				return errors.New("Github api token must be provided.")
			}
			ghp := github.NewProvider(core.Config.Github.Token, DB)
			providers = append(providers, ghp)
		case core.CgitProviderName:
			log15.Info("Creating cgit provider")
			if core.Config.Bing.Key == "" {
				return errors.New("Bing search key are mandatory for cgit provider")
			}
			cgp := cgit.NewProvider(core.Config.Bing.Key, DB)
			providers = append(providers, cgp)
		case core.BitbucketProviderName:
			log15.Info("Creating bitbucket provider")
			bbp := bitbucket.NewProvider(DB)
			providers = append(providers, bbp)
		default:
			return fmt.Errorf("Provider '%s' not found. Allowed providers: %v",
				p, allowedProviders)
		}

	}
	log15.Info("Watcher", "time", c.WatcherTime)
	f, err := c.getPersistFunction()
	if err != nil {
		return err
	}
	watcher := core.NewWatcher(providers, f, c.WatcherTime, time.Second*15)
	watcher.Start()

	select {}

	return nil
}
