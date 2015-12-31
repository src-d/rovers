package linkedin

import (
	"os"

	. "gopkg.in/check.v1"
)

func (s *linkedInSuite) TestLinkedIn_NewImporter(c *C) {
	var tests = [...]struct {
		options      LinkedInImporterOptions
		isError      bool
		errorPattern string
		cookie       string
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
			isError:      true,
			errorPattern: ".*env var not set",
			cookie:       "",
		},
		{
			options: LinkedInImporterOptions{
				Mode: "all",
			},
			isError: false,
			cookie:  "foo",
		},
		{
			options: LinkedInImporterOptions{
				Mode:     "single",
				CodeName: "foo",
			},
			isError: false,
		},
	}

	for idx, tt := range tests {
		if tt.cookie != "" {
			err := os.Setenv(LinkedInCookieEnv, tt.cookie)
			c.Assert(err, IsNil)
		}

		imp, err := NewLinkedInImporter(tt.options)
		if tt.isError {
			c.Assert(err, ErrorMatches, tt.errorPattern,
				Commentf("%d expected %q, got %q", idx, tt.errorPattern, err),
			)
		} else {
			c.Assert(imp, NotNil, Commentf("%d couldn't create LinkedInImporter", idx))
			c.Assert(err, IsNil, Commentf("%d expected nil, got %q", idx, err))
		}
	}
}
