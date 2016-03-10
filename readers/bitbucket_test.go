package readers

import (
	"net/url"
	"os"

	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestBitbucket_GetRepositories(c *C) {
	if os.Getenv("CI_COMMIT") != "" {
		// NOTE: We are not using any auth for Bitbucket so running this on a CI system
		// almost always fails to run because someone else has already exhausted the
		// number of requests
		c.Skip("not running on CI")
	}

	api := NewBitbucketAPI(client.NewClient(true))

	result, err := api.GetRepositories(url.Values{})
	if err != nil {
		c.Skip("Skipped TestBitbucket_GetRepositories because of API rate limits.")
		return
	}
	c.Assert(err, IsNil)
	c.Assert(result.Next.Query().Get("page"), Equals, "2")
	c.Assert(result.Values, HasLen, 10)
	c.Assert(result.Page, Equals, 1)
	c.Assert(result.Values[0].Links.Html.Href, Equals, "https://bitbucket.org/phlogistonjohn/tweakmsg")

	result, err = api.GetRepositories(result.Next.Query())
	if err != nil {
		c.Skip("Skipped TestBitbucket_GetRepositories because of API rate limits.")
		return
	}
	c.Assert(err, IsNil)
	c.Assert(result.Next.Query().Get("page"), Equals, "3")
	c.Assert(result.Values, HasLen, 10)
	c.Assert(result.Page, Equals, 2)
	c.Assert(result.Values[0].Links.Html.Href, Equals, "https://bitbucket.org/bebac/app-template")
}
