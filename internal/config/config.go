package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"

	"github.com/unanet/eve/pkg/artifactory"
	"github.com/unanet/eve/pkg/scm/github"
	"github.com/unanet/eve/pkg/scm/gitlab"
)

var (
	flagConfig *FlagConfig
	config     *Config
	dbConfig   *DBConfig
	mutex      = sync.Mutex{}
)

type LogConfig = log.Config
type ArtifactoryConfig = artifactory.Config
type GitLabConfig = gitlab.Config
type GitHubConfig = github.Config

type DBConfig struct {
	DBHost              string        `envconfig:"DB_HOST" default:"localhost"`
	DBPort              int           `envconfig:"DB_PORT" default:"5432"`
	DBUsername          string        `envconfig:"DB_USERNAME" default:"postgres"`
	DBPassword          string        `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName              string        `envconfig:"DB_NAME" default:"postgres"`
	DBConnectionTimeout time.Duration `envconfig:"DB_CONNECTION_TIMEOUT" default:"10s"`
	LogLevel            string        `envconfig:"LOG_LEVEL" default:"info"`
}

func (c DBConfig) DbConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUsername, c.DBPassword, c.DBName)
}

func (c DBConfig) MigrationConnectionString() string {
	return fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable&user=%s&password=%s", c.DBHost, c.DBPort, c.DBName, c.DBUsername, c.DBPassword)
}

//type IdentityConfig = identity.Config

type Config struct {
	LogConfig
	ArtifactoryConfig
	GitLabConfig
	GitHubConfig
	//Identity               IdentityConfig
	LocalDev               bool          `envconfig:"LOCAL_DEV" default:"false"`
	ApiQUrl                string        `envconfig:"API_Q_URL" required:"true"`
	SourceControlProvider  string        `envconfig:"SCM_PROVIDER" default:"gitlab"`
	ApiQWaitTimeSecond     int64         `envconfig:"API_Q_WAIT_TIME_SECOND" default:"20"`
	ApiQVisibilityTimeout  int64         `envconfig:"API_Q_VISIBILITY_TIMEOUT" default:"3600"`
	ApiQMaxNumberOfMessage int64         `envconfig:"API_Q_MAX_NUMBER_OF_MESSAGE" default:"10"`
	ApiQWorkerTimeout      time.Duration `envconfig:"API_Q_WORKER_TIMEOUT" default:"60s"`
	CronTimeout            time.Duration `envconfig:"CRON_TIMEOUT" default:"120s"`
	HttpCallbackTimeout    time.Duration `envconfig:"HTTP_CALLBACK_TIMEOUT" default:"8s"`
	S3Bucket               string        `envconfig:"S3_BUCKET" required:"true"`
	AWSRegion              string        `envconfig:"AWS_REGION" required:"true"`
	Port                   int           `envconfig:"PORT" default:"8080"`
	MetricsPort            int           `envconfig:"METRICS_PORT" default:"3001"`
	ServiceName            string        `envconfig:"SERVICE_NAME" default:"eve"`
	AdminToken             string        `envconfig:"ADMIN_TOKEN" required:"true"`
}

type FlagConfig struct {
	MigrateFlag bool `split_words:"true" default:"false"`
	ServerFlag  bool `split_words:"true" default:"true"`
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

func GetFlagsConfig() FlagConfig {
	mutex.Lock()
	defer mutex.Unlock()
	if flagConfig != nil {
		return *flagConfig
	}
	c := FlagConfig{}
	err := envconfig.Process("EVE", &c)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
	flagConfig = &c
	return *flagConfig
}
