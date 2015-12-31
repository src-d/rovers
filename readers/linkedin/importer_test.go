package linkedin

import (
	. "gopkg.in/check.v1"
)

func (s *linkedInSuite) TestLinkedIn_NewImporter(c *C) {
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
