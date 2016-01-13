package linkedin

import (
	"gop.kg/src-d/domain@v2.4/models/company"

	. "gopkg.in/check.v1"
)

func (s *linkedInSuite) TestNewLinkedInImporter(c *C) {
	var tests = [...]struct {
		options      LinkedInImporterOptions
		isError      bool
		errorPattern string
	}{
		{
			options: LinkedInImporterOptions{
				Mode:     "all",
				CodeName: "foo",
			},
			isError:      true,
			errorPattern: "supplied codename.*",
		},
		{
			options: LinkedInImporterOptions{
				Mode:     "empty",
				CodeName: "foo",
			},
			isError:      true,
			errorPattern: "supplied codename.*",
		},
		{
			options: LinkedInImporterOptions{
				Mode:     "single",
				CodeName: "",
			},
			isError:      true,
			errorPattern: `.*single requires \-\-codename to be set`,
		},
		{
			options: LinkedInImporterOptions{
				Mode: "foo",
			},
			isError:      true,
			errorPattern: "invalid mode.*",
		},
		{
			options: LinkedInImporterOptions{
				Mode: "all",
			},
			isError: false,
		},
		{
			options: LinkedInImporterOptions{
				Mode:     "single",
				CodeName: "foo",
			},
			isError: false,
		},
		{
			options: LinkedInImporterOptions{
				Mode:   "all",
				Cookie: "",
			},
			isError: false,
		},
		{
			options: LinkedInImporterOptions{
				Mode:   "all",
				Cookie: "eiso",
			},
			isError: false,
		},
	}

	for idx, tt := range tests {
		imp, err := NewLinkedInImporter(tt.options)
		if tt.isError {
			c.Assert(err, ErrorMatches, tt.errorPattern,
				Commentf("%d expected %q, got %q", idx, tt.errorPattern, err),
			)
		} else {
			c.Assert(imp, NotNil)
			c.Assert(err, IsNil)
		}
	}
}

func (s *linkedInSuite) TestNewLinkedInImporter_Save(c *C) {
	var tests = []struct {
		CodeName       string
		Employees      []company.Employee
		AssocEmployees []company.Employee
		IsError        bool
	}{
		{
			CodeName:       "foo",
			Employees:      []company.Employee{employee("John")},
			AssocEmployees: []company.Employee{employee("Mary")},
			IsError:        false,
		},
		{
			CodeName:       "not-exists",
			Employees:      nil,
			AssocEmployees: nil,
			IsError:        true,
		},
	}

	for _, tt := range tests {
		imp, err := NewLinkedInImporter(LinkedInImporterOptions{
			Mode:     "single",
			CodeName: tt.CodeName,
		})
		c.Assert(err, IsNil)

		imp.companyStore = s.store

		err = imp.saveCompanyEmployees(tt.CodeName, tt.Employees, tt.AssocEmployees)
		if tt.IsError {
			c.Assert(err, NotNil)
		} else {
			c.Assert(err, IsNil)

			query := s.store.Query()
			query.FindByCodeName(tt.CodeName)

			doc, err := s.store.FindOne(query)
			c.Assert(err, IsNil)

			c.Assert(doc.Employees, DeepEquals, tt.Employees)
			c.Assert(doc.AssociateEmployees, DeepEquals, tt.AssocEmployees)
		}
	}
}

func employee(name string) company.Employee {
	return company.Employee{
		FirstName:  name,
		LastName:   "Foo",
		Position:   "fooer",
		LinkedInId: 0,
		Location:   "Fooland",
	}
}
