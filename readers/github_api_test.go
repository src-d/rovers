package readers

import . "gopkg.in/check.v1"

func (s *SourcesSuite) TestGithubAPI_GetAllRepositories(c *C) {
	a := NewGithubAPI()
	repos, resp, err := a.GetAllRepositories(0)
	c.Assert(err, IsNil)
	c.Assert(resp.NextPage, Equals, 367)
	c.Assert(repos, HasLen, 100)
}
