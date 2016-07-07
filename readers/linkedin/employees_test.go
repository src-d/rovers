package linkedin

import (
	"os"

	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

const (
	SourcedCompanyId = 10284920
)

func (s *linkedInSuite) TestNewLinkedInWebCrawler(c *C) {
	if os.Getenv("CI_COMMIT") != "" {
		c.Skip("not running on CI")
	}

	cli := client.NewClient(false)
	wc := NewLinkedInWebCrawler(cli, CookieFixtureEiso)
	employees, err := wc.GetEmployees(SourcedCompanyId)
	c.Assert(err, IsNil)
	// NOTE(toqueteos): I'm not sure if checking this value is *that* useful, of
	// course it's needed for validation but it shouldn't be an exact value
	// because changes on the LinkedIn company page are independent of us and
	// may turn up to be a pain in the ass when testing.
	c.Assert(employees, HasLen, 17)
}
