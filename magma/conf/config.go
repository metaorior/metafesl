package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"

)

var Configuration cfg

type COnfiguration struct {

	verboseLevel string `envconfig:"LOG_LEVEL" default:"DEBUG"`

	httpBind string `envconfig: "httpBind" default: "127.0.0.1:80"`
	httpsBind string `envconfig: "httpsBind" default: "127.0.0.1:443"`
	certPath string `envconfig: "certPath" required: "false"`
	privkeyPath string `envconfig: "privkeyPath" required: "false"`

}

func LoadToMemory() {
	if err := envconfig.Process("", &Config); err != nil {
		logrus.WithError(err).Fatal("Initialize configs")
	}
}
