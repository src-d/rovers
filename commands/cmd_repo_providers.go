package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/kr/beanstalk"
	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/github"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	githubProviderName = "github"
	priorityNormal     = 1024
)

var allowedProviders = []string{githubProviderName}

type CmdRrepoProviders struct {
	CmdBase
	Providers   []string      `short:"p" long:"provider" optional:"no" description:"list of providers to execute."`
	GithubToken string        `short:"" long:"github-token" optional:"no" description:"Github API token"`
	WatcherTime time.Duration `short:"t" long:"watcher-time" optional:"no" default:"1h" description:"Time to try again to get new repos"`
	QueueName   string        `short:"q" long:"queue" optional:"no" default:"repo-urls" description:"beanstalkd queue used to send repo urls"`
	Beanstalk   string        `long:"beanstalk" default:"127.0.0.1:11300" description:"beanstalk url server"`
}

func (c *CmdRrepoProviders) Execute(args []string) error {
	c.ChangeLogLevel()

	providers := []core.RepoProvider{}
	for _, p := range c.Providers {
		switch p {
		case githubProviderName:
			log15.Info("Creating github provider")
			if c.GithubToken == "" {
				return errors.New("Github api token must be provided")
			}
			ghp := github.NewProvider(&github.GithubConfig{GithubToken: c.GithubToken})
			providers = append(providers, ghp)
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

func (c *CmdRrepoProviders) getPersistFunction() (func(string) error, error) {
	host := c.Beanstalk
	log15.Info("Beanstalk", "host", host)
	conn, err := beanstalk.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	queue := core.NewBeanstalkQueue(conn, c.QueueName)

	return func(url string) error {
		_, err := queue.Put([]byte(url), priorityNormal, 0, 0)
		if err != nil {
			log15.Error("Error sending data to queue", "error", err)
		}
		return err
	}, nil
}
