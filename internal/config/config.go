package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/artifactory"
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
	Port       int    `split_words:"true" default:"8080"`
	DBHost     string `split_words:"true" default:"localhost"`
	DBPort     int    `split_words:"true" default:"5432"`
	DBUsername string `split_words:"true" default:"eve-api"`
	DBPassword string `split_words:"true" default:"eve-api"`
	DBName     string `split_words:"true" default:"eve-api"`
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
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
}
