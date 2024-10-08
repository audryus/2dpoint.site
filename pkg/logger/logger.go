package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Log interface {
	Info(string, ...map[string]string)
	Error(error, string)
	Core() *zerolog.Logger
}

type log struct {
	logger zerolog.Logger
}

func New() Log {
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	// Set global log level to Info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Sample log messages with a burst of 5 messages and a 10-second pause
	sampler := &zerolog.BurstSampler{Burst: 5, Period: 10 * time.Second}
	logger = logger.Sample(sampler)

	return log{
		logger,
	}
}

func (l log) Core() *zerolog.Logger {
	return &l.logger
}

func (l log) Info(msg string, args ...map[string]string) {
	i := l.logger.Info()
	if len(args) > 0 {
		arg := args[0]
		for key, value := range arg {
			i.Str(key, value)
		}
	}
	i.Msg(msg)
}

func (l log) Error(err error, msg string) {
	l.logger.Error().Stack().Err(err).Msg(msg)
}
