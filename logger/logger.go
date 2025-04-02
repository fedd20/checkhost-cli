package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var loggerInstance *zerolog.Logger

func GetLoggerInstance(debug bool) *zerolog.Logger {
	if loggerInstance == nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger := zerolog.New(writer).With().Timestamp().Logger()
		if debug {
			logger = logger.Level(zerolog.DebugLevel)
		} else {
			logger = logger.Level(zerolog.ErrorLevel)
		}
		loggerInstance = &logger
	}
	return loggerInstance
}
