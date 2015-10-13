package readers

import (
	"testing"

	. "gopkg.in/check.v1"
)

type SourcesSuite struct{}

var _ = Suite(&SourcesSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
