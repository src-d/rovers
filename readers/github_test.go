package readers

import (
	"sort"

	"github.com/src-d/rovers/client"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestGithub_GetProfileByURL_404(c *C) {
	a := NewGithubWebCrawler(client.NewClient(true))
	_, err := a.GetProfileByURL("https://github.com/foobarqux")
	c.Assert(err, Equals, client.NotFound)
}

func (s *SourcesSuite) TestGithub_GetProfileByURL_Company(c *C) {
	a := NewGithubWebCrawler(client.NewClient(true))
	g, err := a.GetProfileByURL("https://github.com/src-d")
	c.Assert(err, IsNil)
	c.Assert(g.Organization, Equals, true)
	c.Assert(g.Username, Equals, "src-d")
	c.Assert(g.FullName, Equals, "source{d}")
	c.Assert(g.Location, Equals, "Madrid, Spain")
	c.Assert(g.Email, Equals, "")
	c.Assert(g.Url, Equals, "https://github.com/src-d")
	members := []string{"Istar-Eldritch", "OmarMohamedDev", "alcortesm", "curratore", "dripolles", "filiptc", "gsc", "ivanfoo", "jorgeschnura", "klaidliadon", "mcuadros", "mvader", "pavelkarpov", "toqueteos"}
	sort.Strings(g.Members)
	sort.Strings(members)
	c.Assert(g.Members, DeepEquals, members)
}

func (s *SourcesSuite) TestGithub_SearchByEmail(c *C) {
	a := NewGithubWebCrawler(client.NewClient(true))
	g, err := a.GetProfileByURL("https://github.com/mcuadros")
	c.Assert(err, IsNil)
	c.Assert(g.Organization, Equals, false)
	c.Assert(g.Username, Equals, "mcuadros")
	c.Assert(g.FullName, Equals, "Máximo Cuadros")
	c.Assert(g.Location, Equals, "Madrid, Spain")
	c.Assert(g.Email, Equals, "mcuadros@gmail.com")
	c.Assert(g.Description, Not(Equals), "")
	c.Assert(g.JoinDate.Unix(), Equals, int64(1332676111))
	sort.Strings(g.Organizations)
	c.Assert(g.Organizations, HasLen, 3)
	c.Assert(g.Organizations, DeepEquals, []string{"/mongator", "/sourcegraph", "/src-d"})
	c.Assert(g.Repositories, HasLen, 5)

	return

	//Change a lot so is hard to test
	c.Assert(g.Repositories[0].Name, Equals, "dockership")
	c.Assert(g.Repositories[0].Url, Equals, "/mcuadros/dockership")
	c.Assert(g.Repositories[0].Owner, Equals, "mcuadros")
	c.Assert(g.Repositories[0].Stars, Not(Equals), 0)
	c.Assert(g.Contributions, HasLen, 5)
	c.Assert(g.Contributions[0].Name, Equals, "mongofill")
	c.Assert(g.Contributions[0].Url, Equals, "/mongofill/mongofill")
	c.Assert(g.Contributions[0].Owner, Equals, "mongofill")
	c.Assert(g.Contributions[0].Stars, Not(Equals), 0)
	c.Assert(g.Followers, Not(Equals), 0)
	c.Assert(g.Starred, Not(Equals), 0)
	c.Assert(g.Following, Not(Equals), 0)
	c.Assert(g.TotalContributions, Not(Equals), 0)

	g, err = a.GetProfileByURL("https://github.com/philips")
	c.Assert(g.Username, Equals, "philips")
	c.Assert(g.WorksFor, Equals, "CoreOS, Inc")
	c.Assert(g.Url, Equals, "https://github.com/philips")
}
