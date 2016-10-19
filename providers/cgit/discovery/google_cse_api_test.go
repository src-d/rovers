package discovery

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

const (
	envKey = "GOOGLE_SEARCH_KEY"
	envCx  = "GOOGLE_SEARCH_CX"
)

type GoogleCseApiSuite struct {
	key string
	cx  string
}

var _ = Suite(&GoogleCseApiSuite{
	key: os.Getenv(envKey),
	cx:  os.Getenv(envCx),
})

func (s *GoogleCseApiSuite) SetUpTest(c *C) {
}

func (s *GoogleCseApiSuite) TestGoogleCseApi_GetPage(c *C) {
	cse := newGoogleCseApi(s.key, s.cx)
	result, err := cse.GetPage(1)

	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	c.Assert(result.Items, NotNil)
}

func (s *GoogleCseApiSuite) TestGoogleCseApi_GetPage_NotExist(c *C) {
	cse := newGoogleCseApi(s.key, s.cx)
	result, err := cse.GetPage(1000000)

	c.Assert(err, Equals, errPageNotFound)
	c.Assert(result, IsNil)
}

func (s *GoogleCseApiSuite) TestGoogleCseApi_PageExists(c *C) {
	cse := newGoogleCseApi(s.key, s.cx)
	result, err := cse.PageExists(1)

	c.Assert(err, IsNil)
	c.Assert(result, Equals, true)

	_, ok := cse.cachedPages[1]
	c.Assert(ok, Equals, true)

	result, err = cse.PageExists(100000)

	c.Assert(err, IsNil)
	c.Assert(result, Equals, false)

	_, ok = cse.cachedPages[100000]
	c.Assert(ok, Equals, false)
}

func (s *GoogleCseApiSuite) TestGoogleCseApi_Reset(c *C) {
	cse := newGoogleCseApi(s.key, s.cx)
	exist, err := cse.PageExists(1)

	c.Assert(err, IsNil)
	c.Assert(exist, Equals, true)

	_, ok := cse.cachedPages[1]
	c.Assert(ok, Equals, true)

	result, err := cse.GetPage(2)

	c.Assert(err, IsNil)
	c.Assert(result, NotNil)

	_, ok = cse.cachedPages[2]
	c.Assert(ok, Equals, true)
}

func (s *GoogleCseApiSuite) TestGoogleCseApi_BadCredentials(c *C) {
	cse := newGoogleCseApi("BAD_KEY", s.cx)
	result, err := cse.GetPage(1)

	c.Assert(err, NotNil)
	c.Assert(err, Equals, errInvalidKey)
	c.Assert(result, IsNil)
}

func Test(t *testing.T) {
	TestingT(t)
}
