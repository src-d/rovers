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
	Cookie      string `short:"" long:"cookie" description:"session cookie to use" required:"true"`
	UseCache    bool   `short:"" long:"cacheUse" description:"wether or not to use the request cache" default:"true"`
	DeleteCache bool   `short:"" long:"cacheDelete" description:"delete cache before running" default:"true"`
	DryRun      bool   `short:"" long:"dry" description:"show employees found, but don't save them" default:"false"`

	companyStore       *models.CompanyStore
	linkedinWebCrawler *linkedin.LinkedInWebCrawler
}

func (cmd *CmdLinkedIn) Execute(args []string) error {
	start := time.Now()

	cli := client.NewClient(cmd.UseCache)
	if cmd.DeleteCache {
		// TODO(toqueteos): Cache should be a storable model, this would be a
		// nice refactor to have by next week (2015Sep28 ~ 2015Oct02)
		// cli.DeleteCache()
	}
	cmd.companyStore = container.GetDomainModelsCompanyStore()
	cmd.linkedinWebCrawler = linkedin.NewLinkedInWebCrawler(cli, cmd.Cookie)

	companiesInfo := cmd.GetCompaniesLinkedInInfo()
	for _, info := range companiesInfo {
		cmd.GetCompanyEmployees(info)
	}

	log15.Info("Done", "elapsed", time.Since(start))

	return nil
}

type CompanyInfo struct {
	CodeName   string
	CompanyIds []int
}

func (cmd *CmdLinkedIn) GetCompaniesLinkedInInfo() []CompanyInfo {
	q := cmd.companyStore.Query()
	set, err := cmd.companyStore.Find(q)
	if err != nil {
		return nil
	}

	var companiesInfo []CompanyInfo
	set.ForEach(func(company *models.Company) error {
		if len(company.LinkedInCompanyIds) == 0 {
			log15.Warn("No LinkedInCompanyIds", "company", company.CodeName)
			return nil
		}

		info := CompanyInfo{
			CodeName:   company.CodeName,
			CompanyIds: company.LinkedInCompanyIds,
		}
		companiesInfo = append(companiesInfo, info)

		return nil
	})

	return companiesInfo
}

func (cmd *CmdLinkedIn) GetCompanyEmployees(info CompanyInfo) {
	var employees []company.Employee
	for _, companyId := range info.CompanyIds {
		subCompanyEmployees, err := cmd.linkedinWebCrawler.GetEmployees(companyId)
		if err != nil {
			log15.Error("Failed to fetch all employees",
				"company", info.CodeName,
				"employees", len(subCompanyEmployees),
				"error", err,
			)
			continue
		}
		employees = append(employees, subCompanyEmployees...)
	}

	cmd.UpdateCompanyEmployees(info, employees)
}

func (cmd *CmdLinkedIn) UpdateCompanyEmployees(info CompanyInfo, employees []company.Employee) {
	if cmd.DryRun {
		cmd.PrintEmployees(employees)
	} else {
		log15.Info("Updating database employees", "company", info.CodeName)
		err := cmd.DoUpdateCompanyEmployees(info.CodeName, employees)
		if err != nil {
			log15.Error("Couldn't update company employees",
				"company", info.CodeName,
				"error", err,
			)
		}
	}
}

func (cmd *CmdLinkedIn) PrintEmployees(employees []company.Employee) {
	for idx, employee := range employees {
		log15.Info("Employee",
			"idx", idx,
			"data", employee,
		)
	}
}

func (cmd *CmdLinkedIn) DoUpdateCompanyEmployees(codeName string, employees []company.Employee) error {
	q := cmd.companyStore.Query()
	q.FindByCodeName(codeName)

	company, err := cmd.companyStore.FindOne(q)
	if err != nil {
		return err
	}
	company.Employees = employees
	_, err = cmd.companyStore.Save(company)
	return err
}
