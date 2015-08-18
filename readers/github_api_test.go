package readers

import . "gopkg.in/check.v1"

func (s *SourcesSuite) TestGithubAPI_GetAllRepositories(c *C) {
	a := NewGithubAPIReader(nil)
	repos, resp, err := a.GetAllRepositories(0)
	c.Assert(err, IsNil)
	c.Assert(resp.NextPage, Equals, 364)
	c.Assert(repos, HasLen, 100)
}
