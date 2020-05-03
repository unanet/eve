package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/secrets"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/gitlab"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

var (
	values *Config
	mutex  = sync.Mutex{}
)

type LogConfig = log.Config
type ArtifactoryConfig = artifactory.Config
type MuxConfig = mux.Config
type GitlabConfig = gitlab.Config
type VaultConfig = secrets.Config

type Config struct {
	LogConfig
	ArtifactoryConfig
	MuxConfig
	GitlabConfig
	VaultConfig
	ApiQUrl                string        `split_words:"true" required:"true"`
	SchQUrl                string        `split_words:"true" required:"true"`
	ApiQWaitTimeSecond     int64         `split_words:"true" default:"20"`
	ApiQVisibilityTimeout  int64         `split_words:"true" default:"3600"`
	ApiQMaxNumberOfMessage int64         `split_words:"true" default:"10"`
	ApiQWorkerTimeout      time.Duration `split_words:"true" default:"60s"`
	S3Bucket               string        `split_words:"true" required:"true"`
	AWSRegion              string        `split_words:"true" required:"true"`
	DBHost                 string        `split_words:"true" default:"localhost"`
	DBPort                 int           `split_words:"true" default:"5432"`
	DBUsername             string        `split_words:"true" default:"eve-api"`
	DBPassword             string        `split_words:"true" default:"eve-api"`
	DBName                 string        `split_words:"true" default:"eve-api"`
}

func GetConfig() Config {
	mutex.Lock()
	defer mutex.Unlock()
	if values != nil {
		return *values
	}
	c := Config{}
	err := envconfig.Process("EVE", &c)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
	values = &c
	return *values
}

func (c Config) DbConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUsername, c.DBPassword, c.DBName)
}

func (c Config) MigrationConnectionString() string {
	return fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable&user=%s&password=%s", c.DBHost, c.DBPort, c.DBName, c.DBUsername, c.DBPassword)
}