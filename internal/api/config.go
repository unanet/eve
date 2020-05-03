package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/gitlab"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

var (
	config   *Config
	dbConfig *DBConfig
	mutex    = sync.Mutex{}
)

type LogConfig = log.Config
type ArtifactoryConfig = artifactory.Config
type MuxConfig = mux.Config
type GitlabConfig = gitlab.Config

type DBConfig struct {
	DBHost              string        `split_words:"true" default:"localhost"`
	DBPort              int           `split_words:"true" default:"5432"`
	DBUsername          string        `split_words:"true" default:"eve-api"`
	DBPassword          string        `split_words:"true" default:"eve-api"`
	DBName              string        `split_words:"true" default:"eve-api"`
	DBConnectionTimeout time.Duration `split_words:"true" default:"10s"`
	LogLevel            string        `split_words:"true" default:"info"`
}

func (c DBConfig) DbConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUsername, c.DBPassword, c.DBName)
}

func (c DBConfig) MigrationConnectionString() string {
	return fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable&user=%s&password=%s", c.DBHost, c.DBPort, c.DBName, c.DBUsername, c.DBPassword)
}

type Config struct {
	LogConfig
	ArtifactoryConfig
	MuxConfig
	GitlabConfig
	ApiQUrl                string        `split_words:"true" required:"true"`
	ApiQWaitTimeSecond     int64         `split_words:"true" default:"20"`
	ApiQVisibilityTimeout  int64         `split_words:"true" default:"3600"`
	ApiQMaxNumberOfMessage int64         `split_words:"true" default:"10"`
	ApiQWorkerTimeout      time.Duration `split_words:"true" default:"60s"`
	HttpCallbackTimeout    time.Duration `split_words:"true" default:"8s"`
	S3Bucket               string        `split_words:"true" required:"true"`
	AWSRegion              string        `split_words:"true" required:"true"`
}

func GetDBConfig() DBConfig {
	mutex.Lock()
	defer mutex.Unlock()
	if dbConfig != nil {
		return *dbConfig
	}
	c := DBConfig{}
	err := envconfig.Process("EVE", &c)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
	dbConfig = &c
	return *dbConfig
}

func GetConfig() Config {
	mutex.Lock()
	defer mutex.Unlock()
	if config != nil {
		return *config
	}
	c := Config{}
	err := envconfig.Process("EVE", &c)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
	config = &c
	return *config
}
