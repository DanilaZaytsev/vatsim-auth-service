package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.DebugLevel)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}
