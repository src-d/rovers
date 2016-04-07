package commands

import (
	"github.com/src-d/rovers/readers/linkedin"
	"gop.kg/src-d/domain@v5/container"
)

type CmdLinkedInUpdate struct {
	CmdBase

	Mode        string `long:"mode" description:"which companies to update" required:"true"`
	CodeName    string `long:"codename" description:"required for --mode=single"`
	LinkedInId  int    `long:"linkedinid" description:"required for --mode=single"`
	Cookie      string `long:"cookie" description:"session cookie to use"`
	Force       bool   `long:"force" description:"force an update of employees, deletes all previous records"`
	UseCache    bool   `long:"cacheUse" description:"wether or not to use the request cache"`
	DeleteCache bool   `long:"cacheDelete" description:"delete cache before running"`
	DryRun      bool   `long:"dry" description:"show employees found, but don't save them"`
}

func (cmd *CmdLinkedInUpdate) Execute(args []string) error {
	cmd.ChangeLogLevel()

	imp, err := linkedin.NewLinkedInImporter(linkedin.LinkedInImporterOptions{
		Mode:             cmd.Mode,
		CodeName:         cmd.CodeName,
		LinkedInId:       cmd.LinkedInId,
		Cookie:           cmd.Cookie,
		UseCache:         cmd.UseCache,
		DeleteCache:      cmd.DeleteCache,
		DryRun:           cmd.DryRun,
		Force:            cmd.Force,
		CompanyStore:     container.GetDomainModelsCompanyStore(),
		CompanyInfoStore: container.GetDomainModelsCompanyInfoStore(),
	})
	if err != nil {
		return err
	}
	err = imp.Import()

	return err
}
