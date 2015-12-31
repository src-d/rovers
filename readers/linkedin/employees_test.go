package linkedin

import (
	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

const (
	CookieFixture = ""
	TybaCompanyId = 924688
)

func (s *linkedInSuite) TestLinkedIn_GetEmployees(c *C) {
	// NOTE: LinkedIn cookie is set via an environment variable. We test by
	// manually commenting this `Skip` or on production.
	c.Skip("Run this locally to test it works")

	cli := client.NewClient(false)
	wc := NewLinkedInWebCrawler(cli, CookieFixture)
	employees, err := wc.GetEmployees(TybaCompanyId)
	c.Assert(err, IsNil)
	// NOTE(toqueteos): I'm not sure if checking this value is *that* useful, of
	// course it's needed for validation but it shouldn't be an exact value
	// because changes on the LinkedIn company page are independent of us and
	// may turn up to be a pain in the ass when testing.
	c.Assert(len(employees), Equals, 32)
}
