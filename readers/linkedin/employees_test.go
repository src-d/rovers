package linkedin

import (
	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

const (
	TybaCompanyId = 924688
)

func (s *linkedInSuite) TestLinkedIn_GetEmployees(c *C) {
	cli := client.NewClient(false)
	wc := NewLinkedInWebCrawler(cli, CookieFixtureEiso)
	employees, err := wc.GetEmployees(TybaCompanyId)
	c.Assert(err, IsNil)
	// NOTE(toqueteos): I'm not sure if checking this value is *that* useful, of
	// course it's needed for validation but it shouldn't be an exact value
	// because changes on the LinkedIn company page are independent of us and
	// may turn up to be a pain in the ass when testing.
	c.Assert(len(employees), Equals, 32)
}
