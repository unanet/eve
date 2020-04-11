package config

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/artifactory"

	"github.com/kelseyhightower/envconfig"

	"gitlab.unanet.io/devops/eve/pkg/log"
)

var (
	Values Config
)

type LogConfig = log.Config
type ArtifactoryConfig = artifactory.Config

type Config struct {
	LogConfig
	ArtifactoryConfig
	Port       int    `default:"8080"`
	DBHost     string `default:"localhost"`
	DBPort     int    `default:"5432"`
	DBUsername string `default:"eve-api"`
	DBPassword string `default:"eve-api"`
	DBName     string `default:"eve-api"`
}

func (c Config) DbConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUsername, c.DBPassword, c.DBName)
}

func (c Config) MigrationConnectionString() string {
	return fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable&user=%s&password=%s", c.DBHost, c.DBPort, c.DBName, c.DBUsername, c.DBPassword)
}

func init() {
	err := envconfig.Process("EVE", &Values)
	if err != nil {
		log.Logger.WithField("error", err).Panic("Unable to Load Config")
	}
}
