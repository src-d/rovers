package readers

import (
	"testing"

	"github.com/tyba/oss/sources/social/http"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type SourcesSuite struct{}

var _ = Suite(&SourcesSuite{})

func (s *SourcesSuite) TestAugurReader_SearchByEmail(c *C) {
	a := NewAugurReader(http.NewClient(true))
	r, res, err := a.SearchByEmail("nawar.alsafar126@gmail.com")

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(r.Status, Equals, 200)
}

func (s *SourcesSuite) TestAugurReader_SearchByEmailBadRequest(c *C) {
	a := NewAugurReader(http.NewClient(true))
	r, res, err := a.SearchByEmail("foo")

	c.Assert(r, IsNil)
	c.Assert(err, Equals, ErrUnexpectedStatusCode)
	c.Assert(res.StatusCode, Equals, 403)
}
