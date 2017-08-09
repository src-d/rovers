package test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
)

func LoadAsset(baseURL string, assetPath string, c *C) {
	filepath.Walk(assetPath, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		responder := ResponderByFile(c, p)
		url := baseURL + strings.Replace(p, assetPath, "", 1)
		httpmock.RegisterResponder(
			"GET",
			url,
			responder,
		)

		if strings.HasSuffix(p, "index.html") {
			url = baseURL + strings.Replace(path.Dir(p), assetPath, "", 1)

			httpmock.RegisterResponder(
				"GET",
				url,
				responder,
			)

			if !strings.HasSuffix(url, "/") {
				url = url + "/"

				httpmock.RegisterResponder(
					"GET",
					url,
					responder,
				)
			}
		}

		return nil
	})
}

func ResponderByFile(c *C, file string) httpmock.Responder {
	return ResponderByFileAndStatus(c, file, 200)
}

func ResponderByFileAndStatus(c *C, file string, status int) httpmock.Responder {
	f, err := os.Open(file)
	c.Assert(err, IsNil)
	data, err := ioutil.ReadAll(f)
	c.Assert(err, IsNil)

	res := httpmock.NewBytesResponse(status, data)
	return httpmock.ResponderFromResponse(res)
}
