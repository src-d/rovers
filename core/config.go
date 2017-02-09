package core

import (
	"github.com/mcuadros/go-defaults"
	"github.com/src-d/envconfig"
)

var (
	Config *config = &config{}
)

type config struct {
	Bing struct {
		Key string
	}
	Github struct {
		Token string
	}
	Postgres struct {
		Url string `default:"postgres://postgres:mysecretpassword@0.0.0.0:5432/postgres?sslmode=disable"`
	}
}

func init() {
	envconfig.MustProcess("config", Config)
	defaults.SetDefaults(Config)
}
