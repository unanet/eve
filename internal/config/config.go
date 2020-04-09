package config

import (
	"github.com/kelseyhightower/envconfig"

	"gitlab.unanet.io/devops/eve/internal/log"
)

var (
	Values Config
)

type Config struct {
	log.Config
	Port       int    `default:"8080"`
	DBHost     string `required:"true"`
	DBUsername string `required:"true"`
	DBPassword string `required:"true"`
	DBName     string `default:"eve-api"`
}

func init() {
	err := envconfig.Process("EVE", &Values)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", err)
		panic(err)
	}
}
