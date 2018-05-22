package logger

import (
	"os"

	"github.com/rs/zerolog"

	"example.com/test/pkg/env"
)

// NewLogger returns a logger
func NewLogger(e *env.Env) zerolog.Logger {

	InitLogger(e)

	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if e.Get("APP_PRETTY_LOGS") != "" && e.Get("APP_PRETTY_LOGS") != "0" {
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return l
}

// InitLogger initializes logger
func InitLogger(e *env.Env) {

	// log level
	logLevel := e.Get("APP_LOG_LEVEL")

	// Default log level is error
	zeroLogLevel := zerolog.ErrorLevel

	if logLevel == "debug" {
		zeroLogLevel = zerolog.DebugLevel
	}

	if logLevel == "info" {
		zeroLogLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zeroLogLevel)
}
