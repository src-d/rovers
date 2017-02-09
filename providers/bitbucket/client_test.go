package bitbucket

import (
	. "gopkg.in/check.v1"
)

const secondPage = "2008-08-13T20:24:13.972049+00:00"

type ClientSuite struct {
	c *client
}

var _ = Suite(&ClientSuite{})

func (s *ClientSuite) SetUpTest(c *C) {
	s.c = newClient()
}

func (s *ClientSuite) TestClient_Repositories(c *C) {
	resp, err := s.c.Repositories("")
	c.Assert(err, IsNil)
	c.Assert(len(resp.Repositories), Equals, pagelenValue)
	c.Assert(resp.Next, Equals, secondPage)
}

func (s *ClientSuite) TestClient_Repositories_Next(c *C) {
	resp, err := s.c.Repositories(secondPage)
	c.Assert(err, IsNil)

	c.Assert(len(resp.Repositories), Equals, pagelenValue)
	c.Assert(resp.Next, Not(Equals), secondPage)

}
