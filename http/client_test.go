package http

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ClientSuite struct{}

var _ = Suite(&ClientSuite{})

func (s *ClientSuite) TestDoHTML(c *C) {
	req, _ := NewRequest("http://httpbin.org/")

	cli := NewClient(true)
	doc, res, err := cli.DoHTML(req)

	c.Assert(doc.Has("body").Nodes, HasLen, 1)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(err, IsNil)
}

type JSONResponse struct {
	Headers interface{}
	Gzipped bool
}

func (s *ClientSuite) TestDoJSON(c *C) {
	req, _ := NewRequest("http://httpbin.org/headers")

	result := &JSONResponse{}

	cli := NewClient(true)
	res, err := cli.DoJSON(req, result)
	c.Assert(len(result.Headers.(map[string]interface{})), Equals, 5)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(err, IsNil)
}

func (s *ClientSuite) TestDoJSONOver400(c *C) {
	req, _ := NewRequest("http://httpbin.org/status/418")

	result := &JSONResponse{}

	cli := NewClient(true)
	res, err := cli.DoJSON(req, result)
	c.Assert(res.StatusCode, Equals, 418)
	c.Assert(err, IsNil)
}

func (s *ClientSuite) TestDoJSONGzip(c *C) {
	req, _ := NewRequest("http://httpbin.org/gzip")

	result := &JSONResponse{}

	cli := NewClient(true)
	res, err := cli.DoJSON(req, result)
	c.Assert(result.Gzipped, Equals, true)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(err, IsNil)
}
