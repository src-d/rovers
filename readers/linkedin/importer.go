package linkedin

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/src-d/rovers/client"
	"gop.kg/src-d/domain@v2.1/container"
	"gop.kg/src-d/domain@v2.1/models"
	"gop.kg/src-d/domain@v2.1/models/company"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2/bson"
)

const LinkedInCookieEnv = "LINKEDIN_COOKIE"

var (
	ErrNoCodeName      = errors.New("--mode=single requires --codename to be set")
	ErrCookieEnvNotSet = fmt.Errorf("%q env var not set", LinkedInCookieEnv)
)

type LinkedInImporterOptions struct {
	Mode        string
	CodeName    string
	UseCache    bool
	DeleteCache bool
	DryRun      bool
}

type LinkedInImporter struct {
	query              bson.M
	options            LinkedInImporterOptions
	companyStore       *models.CompanyStore
	linkedinWebCrawler *LinkedInWebCrawler
}

func NewLinkedInImporter(options LinkedInImporterOptions) (*LinkedInImporter, error) {
	var query bson.M
	switch options.Mode {
	case "all":
		if options.CodeName != "" {
			return nil, fmt.Errorf("supplied codename with --mode=%q", options.Mode)
		}
		query = bson.M{}
	case "empty":
		if options.CodeName != "" {
			return nil, fmt.Errorf("supplied codename with --mode=%q", options.Mode)
		}
		query = bson.M{"employees": bson.M{"$size": 0}}
	case "single":
		if options.CodeName == "" {
			return nil, ErrNoCodeName
		}
		query = bson.M{"codename": options.CodeName}
	default:
		return nil, fmt.Errorf("invalid mode %q", options.Mode)
	}

	cookie := os.Getenv(LinkedInCookieEnv)
	if cookie == "" {
		return nil, ErrCookieEnvNotSet
	}

	cli := client.NewClient(false)
	// if imp.DeleteCache {
	// 	// TODO: Cache as a storable model?
	// 	deleteCache()
	// }
	return &LinkedInImporter{
		query:              query,
		options:            options,
		companyStore:       container.GetDomainModelsCompanyStore(),
		linkedinWebCrawler: NewLinkedInWebCrawler(cli, cookie),
	}, nil
}

func (imp *LinkedInImporter) Import() error {
	start := time.Now()

	companiesInfo := imp.getCompaniesLinkedInInfo()
	for _, info := range companiesInfo {
		companyEmployees := imp.getCompanyEmployees(info)
		associateCompanyEmployees := imp.getAssociateCompanyEmployees(info)

		err := imp.updateCompanyEmployees(info, companyEmployees, associateCompanyEmployees)
		if err != nil {
			log15.Error("Failed to update company employees", "error", err)
			return err
		}
	}

	log15.Info("Done", "elapsed", time.Since(start))
	return nil
}

type CompanyInfo struct {
	CodeName            string
	CompanyIds          []int
	AssociateCompanyIds []int
}

func (imp *LinkedInImporter) getCompaniesLinkedInInfo() []CompanyInfo {
	q := imp.companyStore.Query()
	q.AddCriteria(imp.query)
	set, err := imp.companyStore.Find(q)
	if err != nil {
		return nil
	}

	var companiesInfo []CompanyInfo
	set.ForEach(func(company *models.Company) error {
		if len(company.LinkedInCompanyIds) == 0 && len(company.AssociateCompanyIds) == 0 {
			log15.Warn("No company IDs", "company", company.CodeName)
			return nil
		}

		info := CompanyInfo{
			CodeName:            company.CodeName,
			CompanyIds:          company.LinkedInCompanyIds,
			AssociateCompanyIds: company.AssociateCompanyIds,
		}
		companiesInfo = append(companiesInfo, info)

		return nil
	})

	return companiesInfo
}

func (imp *LinkedInImporter) getCompanyEmployees(info CompanyInfo) []company.Employee {
	return imp.getEmployees(info.CodeName, info.CompanyIds)
}

func (imp *LinkedInImporter) getAssociateCompanyEmployees(info CompanyInfo) []company.Employee {
	return imp.getEmployees(info.CodeName, info.AssociateCompanyIds)
}

func (imp *LinkedInImporter) getEmployees(codeName string, ids []int) []company.Employee {
	var employees []company.Employee
	for _, companyId := range ids {
		companyEmployees, err := imp.linkedinWebCrawler.GetEmployees(companyId)
		if err != nil {
			log15.Error("Failed to fetch all employees",
				"company", codeName,
				"employees", len(companyEmployees),
				"error", err,
			)
			continue
		}
		employees = append(employees, companyEmployees...)
	}

	return employees
}

func (imp *LinkedInImporter) updateCompanyEmployees(
	info CompanyInfo,
	employees, associateEmployees []company.Employee,
) error {
	if imp.options.DryRun {
		log15.Debug("Company employees")
		imp.printEmployees(employees)
		log15.Debug("Associate company employees")
		imp.printEmployees(associateEmployees)
	} else {
		log15.Info("Updating database employees",
			"company", info.CodeName,
			"company_employees", len(employees),
			"associate_employees", len(associateEmployees),
		)
		err := imp.saveCompanyEmployees(info.CodeName, employees, associateEmployees)
		if err != nil {
			log15.Error("Couldn't update company employees",
				"company", info.CodeName,
				"error", err,
			)
			return err
		}
	}
	return nil
}

func (imp *LinkedInImporter) printEmployees(employees []company.Employee) {
	for idx, employee := range employees {
		log15.Debug("Employee",
			"idx", idx,
			"data", employee,
		)
	}
}

func (imp *LinkedInImporter) saveCompanyEmployees(
	codeName string,
	employees, associateEmployees []company.Employee,
) error {
	q := imp.companyStore.Query()
	q.FindByCodeName(codeName)

	company, err := imp.companyStore.FindOne(q)
	if err != nil {
		return err
	}

	var save = false
	if len(employees) > len(company.Employees) {
		log15.Warn("Found less employees",
			"before", len(employees),
			"now", len(company.Employees),
		)
	} else {
		company.Employees = employees
		save = true
	}
	if len(associateEmployees) > len(company.AssociateEmployees) {
		log15.Warn("Found less associate employees",
			"before", len(associateEmployees),
			"now", len(company.AssociateEmployees),
		)
	} else {
		company.AssociateEmployees = associateEmployees
		save = true
	}

	if save {
		_, err = imp.companyStore.Save(company)
		return err
	}

	return nil
}
