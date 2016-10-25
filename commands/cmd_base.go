package commands

import (
	"fmt"
	"os"

	"gopkg.in/inconshreveable/log15.v2"
)

const (
	envGoogleKey   = "GOOGLE_SEARCH_KEY"
	envGoogleCx    = "GOOGLE_SEARCH_CX"
	envGithubToken = "GITHUB_TOKEN"
)

type EnvVars struct {
	googleCSEApiKey string
	googleCSECxKey  string
	githubToken     string
}

type CmdBase struct {
	LogLevel string `short:"" long:"loglevel" description:"max log level enabled" default:"info"`
	LogFile  string `short:"" long:"logfile" description:"path to file where logs will be stored" default:""`
	EnvVars  *EnvVars
}

func (c *CmdBase) InitVars() {
	c.EnvVars = &EnvVars{
		googleCSEApiKey: os.Getenv(envGoogleKey),
		googleCSECxKey:  os.Getenv(envGoogleCx),
		githubToken:     os.Getenv(envGithubToken),
	}
}

func (c *CmdBase) ChangeLogLevel() {
	lvl, err := log15.LvlFromString(c.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("unknown level name %q", c.LogLevel))
	}

	handler := log15.StdoutHandler
	if c.LogFile != "" {
		handler = log15.MultiHandler(
			log15.StdoutHandler,
			log15.Must.FileHandler(c.LogFile, log15.LogfmtFormat()),
		)
	}
	log15.Root().SetHandler(log15.LvlFilterHandler(lvl, handler))
}
