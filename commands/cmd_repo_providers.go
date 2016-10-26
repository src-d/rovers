package commands

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/kr/beanstalk"
	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit"
	"github.com/src-d/rovers/providers/github"
	"gop.kg/src-d/domain@v6/models/repository"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	githubProviderName = "github"
	cgitProviderName   = "cgit"
	priorityNormal     = 1024
)

var allowedProviders = []string{githubProviderName, cgitProviderName}

type CmdRepoProviders struct {
	CmdBase
	Providers   []string      `short:"p" long:"provider" optional:"no" description:"list of providers to execute."`
	WatcherTime time.Duration `short:"t" long:"watcher-time" optional:"no" default:"1h" description:"Time to try again to get new repos"`
	QueueName   string        `short:"q" long:"queue" optional:"no" default:"repo-urls" description:"beanstalkd queue used to send repo urls"`
	Beanstalk   string        `long:"beanstalk" default:"127.0.0.1:11300" description:"beanstalk url server"`
}

func (c *CmdRepoProviders) Execute(args []string) error {
	c.InitVars()
	c.ChangeLogLevel()

	providers := []core.RepoProvider{}
	for _, p := range c.Providers {
		switch p {
		case githubProviderName:
			log15.Info("Creating github provider")
			if c.EnvVars.githubToken == "" {
				return fmt.Errorf("Github api token must be provided. Env variable: %s", envGithubToken)
			}
			ghp := github.NewProvider(&github.GithubConfig{GithubToken: c.EnvVars.githubToken})
			providers = append(providers, ghp)
		case cgitProviderName:
			log15.Info("Creating cgit provider")
			if c.EnvVars.googleCSECxKey == "" || c.EnvVars.googleCSEApiKey == "" {
				return fmt.Errorf("Environment variables %s and %s are mandatory for cgit provider",
					envGoogleKey, envGoogleCx)
			}
			cgp := cgit.NewProvider(c.EnvVars.googleCSEApiKey, c.EnvVars.googleCSECxKey)
			providers = append(providers, cgp)
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
	host := c.Beanstalk
	log15.Info("Beanstalk", "host", host)
	conn, err := beanstalk.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	queue := core.NewBeanstalkQueue(conn, c.QueueName)

	return func(repo *repository.Raw) error {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(repo)
		if err != nil {
			log15.Error("gob.Encode", "error", err)
			return err
		}
		_, err = queue.Put(buf.Bytes(), priorityNormal, 0, 0)
		if err != nil {
			log15.Error("Error sending data to queue", "error", err)
			return err
		}
		return nil
	}, nil
}
