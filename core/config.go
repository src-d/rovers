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
	Broker struct {
		URL string `default:"amqp://guest:guest@localhost:5672/"`
	}
}

func init() {
	envconfig.MustProcess("config", Config)
	defaults.SetDefaults(Config)
}
