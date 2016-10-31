package core

import (
	"github.com/mcuadros/go-defaults"
	"github.com/src-d/envconfig"
)

var (
	Config *config = &config{}
)

type config struct {
	Google struct {
		SearchKey string
		SearchCx  string
	}
	Github struct {
		Token string
	}
	MongoDb struct {
		Database struct {
			Github string `default:"github"`
			Cgit   string `default:"cgit"`
		}
	}
}

func init() {
	envconfig.MustProcess("config", Config)
	defaults.SetDefaults(Config)
}
