package linkedin

import (
	"errors"
	"fmt"
	"time"

	"github.com/src-d/rovers/client"
	"gop.kg/src-d/domain@v5/models"
	"gop.kg/src-d/domain@v5/models/company"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/src-d/storable.v1"
)

var (
	ErrBadArguments = errors.New("no LinkedIn Ids provided")
	ErrSingleParam  = errors.New("--mode=single requires --codename or --linkedinid to be set")
	ErrNoCodeName   = errors.New("--mode=single requires --codename to be set")
	ErrNoLinkedInId = errors.New("--mode=single requires --linkedinid to be set")
	ErrNoCookie     = errors.New("--cookie not set")
)

type LinkedInImporterOptions struct {
	Mode        string
	CodeName    string
	LinkedInId  int
	Cookie      string
	UseCache    bool
	DeleteCache bool
	DryRun      bool
	Force       bool

	CompanyStore     *models.CompanyStore
	CompanyInfoStore *models.CompanyInfoStore
}

type LinkedInImporter struct {
	ids                []int
	options            LinkedInImporterOptions
	linkedinWebCrawler *LinkedInWebCrawler
}

func NewLinkedInImporter(options LinkedInImporterOptions) (*LinkedInImporter, error) {
	var ids []int

	switch options.Mode {
	case "all":
		if options.CodeName != "" {
			return nil, fmt.Errorf("supplied codename with --mode=%q", options.Mode)
		}
		if options.LinkedInId != 0 {
			return nil, fmt.Errorf("supplied linkedinid with --mode=%q", options.Mode)
		}

		ids = getIdsFromQuery(options.CompanyInfoStore, bson.M{})
	case "empty":
		if options.CodeName != "" {
			return nil, fmt.Errorf("supplied codename with --mode=%q", options.Mode)
		}
		if options.LinkedInId != 0 {
			return nil, fmt.Errorf("supplied linkedinid with --mode=%q", options.Mode)
		}

		ids = getIdsFromQuery(options.CompanyInfoStore, bson.M{"employees": bson.M{"$size": 0}})
	case "single":
		if options.CodeName == "" && options.LinkedInId == 0 {
			return nil, ErrSingleParam
		}

		if options.CodeName != "" {
			ids = getIdsFromCodeName(options.CompanyStore, options.CodeName)
		} else if options.LinkedInId != 0 {
			ids = append(ids, options.LinkedInId)
		}
	default:
		return nil, fmt.Errorf("invalid mode %q", options.Mode)
	}

	// Future-proof: we may end up using Jorge's or Ivan's cookie too
	// --cookie is not required anymore
	switch options.Cookie {
	case "", "eiso":
		options.Cookie = CookieFixtureEiso
	}

	cli := client.NewClient(options.UseCache)
	// if imp.DeleteCache {
	// 	// TODO: Cache as a storable model?
	// 	deleteCache()
	// }

	return &LinkedInImporter{
		ids:                ids,
		options:            options,
		linkedinWebCrawler: NewLinkedInWebCrawler(cli, options.Cookie),
	}, nil
}

func (imp *LinkedInImporter) Import() error {
	start := time.Now()

	if imp.ids == nil || len(imp.ids) == 0 {
		return ErrBadArguments
	}

	for _, id := range imp.ids {
		employees, err := imp.getEmployees(id)
		if err != nil {
			log15.Error("Failed to get company employees",
				"id", id,
				"error", err,
			)

			continue
		}

		if err := imp.saveEmployees(id, employees); err != nil {
			log15.Error("Failed to save company employees",
				"id", id,
				"employees", len(employees),
				"error", err,
			)

			continue
		}
	}

	log15.Info("Import done", "elapsed", time.Since(start))

	return nil
}

func (imp *LinkedInImporter) getEmployees(linkedInId int) ([]company.Employee, error) {
	return imp.linkedinWebCrawler.GetEmployees(linkedInId)
}

func (imp *LinkedInImporter) saveEmployees(
	linkedInId int,
	employees []company.Employee,
) error {
	log15.Info("Saving employees",
		"id", linkedInId,
		"employees", len(employees),
	)

	if imp.options.DryRun {
		log15.Warn("--dry supplied, not actually saving")
		return nil
	}

	return imp.save(linkedInId, employees)
}

func (imp *LinkedInImporter) save(linkedInId int, employees []company.Employee) error {
	query := imp.options.CompanyInfoStore.Query()
	query.FindByLinkedInId(linkedInId)

	info, err := imp.options.CompanyInfoStore.FindOne(query)
	if err == storable.ErrNotFound {
		info = imp.options.CompanyInfoStore.New(linkedInId)
	} else if err != nil {
		return err
	}

	oldNumber := len(info.Employees)
	newNumber := len(employees)
	if imp.options.Force {
		log15.Info("--force was provided, ignoring safeguards")
	} else {
		if newNumber < (oldNumber / 2) {
			log15.Crit("Safeguard triggered",
				"id", linkedInId,
				"old", oldNumber,
				"new", newNumber,
			)

			return fmt.Errorf("Found %d employees for company #%d, it had %d",
				newNumber, linkedInId, oldNumber)
		}
	}

	info.Employees = employees
	_, err = imp.options.CompanyInfoStore.Save(info)

	return err
}

func getIdsFromQuery(
	store *models.CompanyInfoStore,
	criteria bson.M,
) (ids []int) {
	query := store.Query()
	query.AddCriteria(criteria)

	set, err := store.Find(query)
	if err != nil {
		log15.Error("Can't find companies",
			"query", criteria,
			"error", err,
		)
		return nil
	}

	err = set.ForEach(func(doc *models.CompanyInfo) error {
		ids = append(ids, doc.LinkedInId)

		return nil
	})
	if err != nil {
		log15.Error("Error iterating companies",
			"query", criteria,
			"error", err,
		)
		return nil
	}

	return ids
}

func getIdsFromCodeName(
	store *models.CompanyStore,
	codeName string,
) (ids []int) {
	query := store.Query()
	query.FindByCodeName(codeName)

	set, err := store.Find(query)
	if err != nil {
		log15.Error("Can't find companies",
			"codename", codeName,
			"error", err,
		)
		return nil
	}

	err = set.ForEach(func(doc *models.Company) error {
		ids = append(ids, doc.LinkedInCompanyIds...)
		ids = append(ids, doc.AssociateCompanyIds...)

		return nil
	})
	if err != nil {
		log15.Error("Error iterating companies",
			"codename", codeName,
			"error", err,
		)
		return nil
	}

	return ids
}
