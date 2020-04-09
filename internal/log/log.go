package log

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	logLevel string `default:"info"`
}

var (
	Logger *logrus.Logger
)

func init() {
	var c Config
	configErr := envconfig.Process("EVE", &c)
	// Provide default values for the Logger in case it can't load the config and then
	// log the error after the logger is loaded.

	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	level, err := logrus.ParseLevel(strings.ToLower(c.logLevel))
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
	if configErr != nil {
		Logger.Error("Logger Config failed to Load", configErr)
	}
}
