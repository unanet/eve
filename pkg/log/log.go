package log

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel string `default:"info"`
}

var (
	Logger     *logrus.Logger
	httpLogger *logrus.Logger
)

func init() {
	var c Config
	configErr := envconfig.Process("EVE", &c)
	// Provide default values for the Logger in case it can't load the config and then
	// log the error after the logger is loaded.

	Logger = logrus.New()
	httpLogger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	httpLogger.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
	level, err := logrus.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
	httpLogger.SetLevel(level)
	if configErr != nil {
		Logger.Error("Logger Config failed to Load", configErr)
	}
}
