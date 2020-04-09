package config

import (
	"github.com/kelseyhightower/envconfig"

	"gitlab.unanet.io/devops/eve/pkg/log"
)

var (
	Values Config
)

type Config struct {
	log.Config
	Port       int    `default:"8080"`
	DBHost     string `default:"localhost"`
	DBPort     string `default:"5432"`
	DBUsername string `default:"eve-api"`
	DBPassword string `default:"eve-api"`
	DBName     string `default:"eve-api"`
}

func init() {
	err := envconfig.Process("EVE", &Values)
	if err != nil {
		log.Logger.WithField("error", err).Panic("Unable to Load Config")
	}
}
