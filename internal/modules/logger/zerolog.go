package logger

import (
	"time"

	"github.com/rs/zerolog"
)

var _ Logger = (*zerologLogger)(nil)

type zerologLogger struct {
	config *config
	logger *zerolog.Logger
}

func NewZerolog(c *config) *zerologLogger {
	writer := zerolog.ConsoleWriter{
		Out:        c.Writer,
		NoColor:    true,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(writer).
		Level(zerolog.Level(c.Level)).
		With().
		Timestamp().
		Caller().
		Logger()

	return &zerologLogger{
		config: c,
		logger: &logger,
	}
}

func (l *zerologLogger) Config() config {
	return *l.config
}

func (l *zerologLogger) D(format string, args ...interface{}) {
	l.logger.Debug().CallerSkipFrame(1).Msgf(format, args...)
}

func (l *zerologLogger) I(format string, args ...interface{}) {
	l.logger.Info().CallerSkipFrame(1).Msgf(format, args...)
}

func (l *zerologLogger) W(format string, args ...interface{}) {
	l.logger.Warn().CallerSkipFrame(1).Msgf(format, args...)
}

func (l *zerologLogger) E(format string, args ...interface{}) {
	l.logger.Error().CallerSkipFrame(1).Msgf(format, args...)
}

func (l *zerologLogger) C(format string, args ...interface{}) {
	l.logger.Fatal().CallerSkipFrame(1).Msgf(format, args...)
}

func (l *zerologLogger) Fields(fields Fields) Logger {
	newLogger := l.logger.With()
	for k, v := range fields {
		newLogger = newLogger.Interface(k, v)
	}
	logger := newLogger.Logger()

	return &zerologLogger{
		config: l.config,
		logger: &logger,
	}
}

func (l *zerologLogger) Err(err error) Logger {
	logger := l.logger.With().Err(err).Logger()

	return &zerologLogger{
		config: l.config,
		logger: &logger,
	}
}
