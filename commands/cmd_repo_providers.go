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
	ocore "srcd.works/core.v0"
	"srcd.works/core.v0/model"
	"srcd.works/framework.v0/queue"
)

const (
	githubProviderName    = "github"
	cgitProviderName      = "cgit"
	bitbucketProviderName = "bitbucket"

	priorityNormal = 1024
)

var allowedProviders = []string{githubProviderName, cgitProviderName, bitbucketProviderName}

type CmdRepoProviders struct {
	CmdBase
	Providers   []string      `short:"p" long:"provider" optional:"yes" description:"list of providers to execute. (default: all)"`
	WatcherTime time.Duration `short:"t" long:"watcher-time" optional:"no" default:"1h" description:"Time to try again to get new repos"`
	Queue       string        `long:"queue" default:"rovers" description:"queue name"`
}

func (c *CmdRepoProviders) Execute(args []string) error {
	c.ChangeLogLevel()

	if len(c.Providers) == 0 {
		log15.Info("No providers added using --provider option. Executing all known providers",
			"providers", allowedProviders)
		c.Providers = allowedProviders
	}

	DB := ocore.Database()

	providers := []core.RepoProvider{}
	for _, p := range c.Providers {
		switch p {
		case githubProviderName:
			log15.Info("Creating github provider")
			if core.Config.Github.Token == "" {
				return errors.New("Github api token must be provided.")
			}
			ghp := github.NewProvider(core.Config.Github.Token, DB)
			providers = append(providers, ghp)
		case cgitProviderName:
			log15.Info("Creating cgit provider")
			if core.Config.Bing.Key == "" {
				return errors.New("Bing search key are mandatory for cgit provider")
			}
			cgp := cgit.NewProvider(core.Config.Bing.Key, DB)
			providers = append(providers, cgp)
		case bitbucketProviderName:
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
	return nil
}

func (c *CmdRepoProviders) getPersistFunction() (core.PersistFN, error) {
	broker, err := queue.NewBroker(core.Config.Broker.URL)
	if err != nil {
		return nil, err
	}

	q, err := broker.Queue(c.Queue)
	if err != nil {
		return nil, err
	}

	return func(repo *model.Mention) error {
		j := queue.NewJob()

		if err := j.Encode(repo); err != nil {
			return err
		}

		return q.Publish(j)
	}, nil
}
