package commands

import (
	"fmt"

	"gopkg.in/inconshreveable/log15.v2"
)

type CmdBase struct {
	LogLevel string `short:"" long:"loglevel" description:"max log level enabled" default:"info"`
	LogFile  string `short:"" long:"logfile" description:"path to file where logs will be stored" default:""`
}

func (c *CmdBase) ChangeLogLevel() {
	lvl, err := log15.LvlFromString(c.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("unknown level name %q", c.LogLevel))
	}

	handlers := []log15.Handler{log15.CallerFileHandler(log15.StdoutHandler)}
	if c.LogFile != "" {
		handlers = append(handlers,
			log15.CallerFileHandler(log15.Must.FileHandler(c.LogFile, log15.LogfmtFormat())))
	}
	log15.Root().SetHandler(log15.LvlFilterHandler(lvl, log15.MultiHandler(handlers...)))
}
