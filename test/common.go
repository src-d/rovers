package test

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
)

func LoadAsset(baseURL string, assetPath string, c *C) {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(assetPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		responder := ResponderByFile(c, p)
		p, err = filepath.Rel(assetPath, p)
		if err != nil {
			return err
		}

		//url := baseURL + filepath.ToSlash(p)
		registerResponder(u, filepath.ToSlash(p), responder, false)

		if strings.HasSuffix(p, "index.html") {
			p = filepath.ToSlash(filepath.Dir(p))
			registerResponder(u, p, responder, true)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func registerResponder(u *url.URL, p string, resp httpmock.Responder, index bool) {
	baseURL := *u
	baseURL.Path = path.Join("/", baseURL.Path, p)
	httpmock.RegisterResponder(
		"GET",
		baseURL.String(),
		resp,
	)

	if index && !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path = baseURL.Path + "/"
		httpmock.RegisterResponder(
			"GET",
			baseURL.String(),
			resp,
		)
	}
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
