package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
	"os"
)

func NewLogger() zerolog.Logger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return log.Logger
}

func NewEventLogger(log zerolog.Logger) fxevent.Logger {
	return &fxevent.ConsoleLogger{W: log}
}
