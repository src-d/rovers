package readers

import (
	"net/url"

	"github.com/tyba/srcd-rovers/client"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestBitbucket_GetRepositories(c *C) {
	api := NewBitbucketAPI(client.NewClient(true))

	result, err := api.GetRepositories(url.Values{})
	c.Assert(err, IsNil)
	c.Assert(result.Next.Query().Get("page"), Equals, "2")
	c.Assert(result.Values, HasLen, 10)
	c.Assert(result.Page, Equals, 1)
	c.Assert(result.Values[0].Links.Html.Href, Equals, "https://bitbucket.org/phlogistonjohn/tweakmsg")

	result, err = api.GetRepositories(result.Next.Query())
	c.Assert(err, IsNil)
	c.Assert(result.Next.Query().Get("page"), Equals, "3")
	c.Assert(result.Values, HasLen, 10)
	c.Assert(result.Page, Equals, 2)
	c.Assert(result.Values[0].Links.Html.Href, Equals, "https://bitbucket.org/bebac/app-template")
}
