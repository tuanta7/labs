package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.Logger,
	}
}
