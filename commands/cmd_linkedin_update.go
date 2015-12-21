package commands

import (
	"time"

	"github.com/src-d/rovers/readers/linkedin"

	"gopkg.in/inconshreveable/log15.v2"
)

type CmdLinkedInUpdate struct {
	CmdBase

	Mode        string `long:"mode" description:"which companies to update" required:"true"`
	CodeName    string `long:"codename" description:"required for --mode=single"`
	Cookie      string `long:"cookie" description:"session cookie to use"`
	UseCache    bool   `long:"cacheUse" description:"wether or not to use the request cache" default:"false"`
	DeleteCache bool   `long:"cacheDelete" description:"delete cache before running" default:"false"`
	DryRun      bool   `long:"dry" description:"show employees found, but don't save them" default:"false"`
}

func (cmd *CmdLinkedInUpdate) Execute(args []string) error {
	cmd.ChangeLogLevel()

	start := time.Now()
	imp, err := linkedin.NewLinkedInImporter(linkedin.LinkedInImporterOptions{
		Mode:        cmd.Mode,
		CodeName:    cmd.CodeName,
		Cookie:      cmd.Cookie,
		UseCache:    cmd.UseCache,
		DeleteCache: cmd.DeleteCache,
		DryRun:      cmd.DryRun,
	})
	if err != nil {
		return err
	}
	err = imp.Import()
	log15.Info("Done", "elapsed", time.Since(start))

	return err
}
