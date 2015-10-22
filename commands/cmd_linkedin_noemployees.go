package commands

import (
	"github.com/src-d/domain/models"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2/bson"
)

type CmdLinkedInNoEmployees struct {
	CmdLinkedIn
}

func (cmd *CmdLinkedInNoEmployees) GetCompaniesLinkedInInfo() []CompanyInfo {
	q := cmd.companyStore.Query()
	q.AddCriteria(bson.M{"linkedincompanyids": bson.M{"$size": 0}})
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
