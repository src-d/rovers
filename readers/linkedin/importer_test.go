package linkedin

import (
	"sort"

	"gop.kg/src-d/domain@v5/models/company"

	. "gopkg.in/check.v1"
)

func (s *linkedInSuite) TestNewLinkedInImporter(c *C) {
	var tests = [...]struct {
		options      LinkedInImporterOptions
		ids          []int
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
				Mode:       "all",
				LinkedInId: 1234,
			},
			isError:      true,
			errorPattern: "supplied linkedinid.*",
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
				Mode:       "empty",
				LinkedInId: 1234,
			},
			isError:      true,
			errorPattern: "supplied linkedinid.*",
		},
		{
			options: LinkedInImporterOptions{
				Mode: "empty",
			},
			ids:     []int{1, 2},
			isError: false,
		},
		{
			options: LinkedInImporterOptions{
				Mode:     "single",
				CodeName: "",
			},
			isError:      true,
			errorPattern: ErrSingleParam.Error(),
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
			ids:     []int{1, 2},
			isError: false,
		},
		{
			options: LinkedInImporterOptions{
				Mode:       "single",
				LinkedInId: 1,
			},
			ids:     []int{1},
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
		tt.options.CompanyStore = s.compStore
		tt.options.CompanyInfoStore = s.infoStore

		imp, err := NewLinkedInImporter(tt.options)
		if !tt.isError {
			c.Assert(err, IsNil)
			c.Assert(imp, NotNil,
				Commentf("test case #%d expected NotNil importer", idx))

			if tt.ids != nil {
				sort.Ints(imp.ids)
				sort.Ints(tt.ids)
				c.Assert(imp.ids, DeepEquals, tt.ids,
					Commentf("test case #%d expected ids=%d, got ids=%d", idx, tt.ids, imp.ids))
			}
		} else {
			c.Assert(err, ErrorMatches, tt.errorPattern,
				Commentf("test case #%d expected %q, got %q", idx, tt.errorPattern, err))
			c.Assert(imp, IsNil)
		}
	}
}

func (s *linkedInSuite) TestNewLinkedInImporter_Save(c *C) {
	var tests = []struct {
		LinkedInId int
		CodeName   string
		Employees  []company.Employee
		DryRun     bool
		IsError    bool
	}{
		{
			LinkedInId: 1,
			CodeName:   "foo",
			Employees:  []company.Employee{employee("John One")},
			DryRun:     false,
			IsError:    false,
		},
		{
			LinkedInId: 5,
			CodeName:   "foo",
			Employees:  []company.Employee{employee("John Five")},
			DryRun:     true,
			IsError:    false,
		},
		{
			LinkedInId: 0,
			CodeName:   "not-exists",
			Employees:  []company.Employee{},
			DryRun:     false,
			IsError:    false,
		},
	}

	for idx, tt := range tests {
		imp, err := NewLinkedInImporter(LinkedInImporterOptions{
			Mode:             "single",
			CodeName:         tt.CodeName,
			DryRun:           tt.DryRun,
			CompanyStore:     s.compStore,
			CompanyInfoStore: s.infoStore,
		})
		c.Assert(err, IsNil)

		err = imp.saveEmployees(tt.LinkedInId, tt.Employees)
		if !tt.IsError {
			c.Assert(err, IsNil, Commentf("test case #%d expected a nil error, got %q", idx, err))

			if tt.DryRun {
				continue
			}

			query := s.infoStore.Query()
			query.FindByLinkedInId(tt.LinkedInId)
			doc, err := s.infoStore.FindOne(query)
			c.Assert(err, IsNil)

			c.Assert(doc.LinkedInId, Equals, tt.LinkedInId, Commentf("test case #%d", idx))
			c.Assert(doc.Employees, DeepEquals, tt.Employees, Commentf("test case #%d", idx))
		} else {
			c.Assert(err, NotNil, Commentf("test case #%d expected non nil error, got %q", idx, err))
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
