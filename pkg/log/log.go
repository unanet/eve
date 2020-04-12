package log

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

type Config struct {
	LogLevel       string `split_words:"true" default:"info"`
	LogServiceName string `split_words:"true" default:"eve"`
}

func logLevel(cfgLevel string) zap.AtomicLevel {
	var logLevel zap.AtomicLevel
	switch cfgLevel {
	case "debug":
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error", "err":
		logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		logLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return logLevel
}

func newLogger(sn string, ll string) *zap.Logger {
	cfg := zap.Config{
		Level:            logLevel(strings.ToLower(ll)),
		Encoding:         "json",
		DisableCaller:    true,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"service": sn},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func init() {
	var c Config
	configErr := envconfig.Process("EVE", &c)
	Logger = newLogger(c.LogServiceName, c.LogLevel)

	if configErr != nil {
		Logger.Error("Logger Config failed to Load", zap.Error(configErr))
	}
}
