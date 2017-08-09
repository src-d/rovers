package bitbucket

import (
	"github.com/src-d/rovers/test"
	. "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
)

func LoadAssets(c *C) {
	responder1 := test.ResponderByFile(c, "assets/1.json")
	httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories?pagelen=100", responder1)

	responder2 := test.ResponderByFile(c, "assets/2.json")
	httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories?after=2008-08-13T20%3A43%3A55.039582%2B00%3A00&pagelen=100", responder2)

	responder3 := test.ResponderByFile(c, "assets/3.json")
	httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories?after=2011-08-10T00%3A42%3A35.509559%2B00%3A00&pagelen=100", responder3)

	responder404 := test.ResponderByFileAndStatus(c, "assets/404.json", 404)
	httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories?after=3000-01-00T17%3A25%3A17.038951%2B00%3A00&pagelen=100", responder404)
}
