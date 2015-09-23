package commands

import (
	"time"

	"github.com/tyba/srcd-domain/container"

	"github.com/tyba/srcd-domain/models"
	"github.com/tyba/srcd-domain/models/company"
	"github.com/tyba/srcd-rovers/client"
	"github.com/tyba/srcd-rovers/readers/linkedin"

	"gopkg.in/inconshreveable/log15.v2"
)

type CmdLinkedIn struct {
	CompanyCodeName string `short:"" long:"companyCodename" description:"" required:"true"`
	CompanyId       int    `short:"" long:"companyId" description:"LinkedIn company page Id" required:"true"`
	Cookie          string `short:"" long:"cookie" description:"session cookie to use" required:"true"`
	UseCache        bool   `short:"" long:"cacheUse" description:"wether or not to use the request cache" default:"true"`
	DeleteCache     bool   `short:"" long:"cacheDelete" description:"delete cache before running" default:"true"`
	DryRun          bool   `short:"" long:"dry" description:"show employees found, but don't save them" default:"false"`

	companyStore *models.CompanyStore
}

func (cmd *CmdLinkedIn) Execute(args []string) error {
	cmd.companyStore = container.GetDomainModelsCompanyStore()

	start := time.Now()

	cli := client.NewClient(cmd.UseCache)
	if cmd.DeleteCache {
		// TODO(toqueteos): Cache should be a storable model, this would be a
		// nice refactor to have by next week (2015Sep28 ~ 2015Oct02)

		// cli.DeleteCache()
	}
	wc := linkedin.NewLinkedInWebCrawler(cli, cmd.Cookie)
	employees, err := wc.GetEmployees(cmd.CompanyId)
	if err != nil {
		log15.Error("Failed to fetch all employees",
			"employees", len(employees),
			"error", err,
		)
	}

	log15.Info("Done", "elapsed", time.Since(start))

	if cmd.DryRun {
		cmd.PrintEmployees(employees)
	} else {
		log15.Info("Updating database employees", "company", cmd.CompanyCodeName)
		err = cmd.UpdateCompanyEmployees(employees)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *CmdLinkedIn) PrintEmployees(employees []company.Employee) {
	for idx, employee := range employees {
		log15.Info("Employee",
			"idx", idx,
			"data", employee,
		)
	}
}

func (cmd *CmdLinkedIn) UpdateCompanyEmployees(employees []company.Employee) error {
	q := cmd.companyStore.Query()
	q.FindByCodeName(cmd.CompanyCodeName)

	company, err := cmd.companyStore.FindOne(q)
	if err != nil {
		return err
	}
	company.Employees = employees
	_, err = cmd.companyStore.Save(company)
	return err
}
