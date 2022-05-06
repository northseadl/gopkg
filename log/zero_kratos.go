package log

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"os"
	"time"
)

var _ log.Logger = (*KratosLogger)(nil)
var devLogger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

type KratosLogger struct {
	log *zerolog.Logger
	dev bool
}

func NewKratosLogger(zeroLogger *zerolog.Logger, dev bool) *KratosLogger {
	kLogger := KratosLogger{zeroLogger, dev}
	if dev {
		devLogger.Debug().Msg("zero-dev-logger enabled")
	}
	return &kLogger
}

func (l *KratosLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn().Msg(fmt.Sprint("Key-values must appear in pairs: ", keyvals))
		return nil
	}
	var e *zerolog.Event
	// dev stdout
	if l.dev {
		switch level {
		case log.LevelDebug:
			e = devLogger.Debug()
		case log.LevelInfo:
			e = devLogger.Info()
		case log.LevelWarn:
			e = devLogger.Warn()
		case log.LevelError:
			e = devLogger.Error().Stack()
		case log.LevelFatal:
			e = devLogger.Fatal().Stack()
		}

		for i := 0; i < len(keyvals); i += 2 {
			e = e.Interface(keyvals[i].(string), keyvals[i+1])
		}
		e.Send()
	}
	switch level {
	case log.LevelDebug:
		e = l.log.Debug()
	case log.LevelInfo:
		e = l.log.Info()
	case log.LevelWarn:
		e = l.log.Warn()
	case log.LevelError:
		e = l.log.Error().Stack()
	case log.LevelFatal:
		e = l.log.Fatal().Stack()
	}

	for i := 0; i < len(keyvals); i += 2 {
		e = e.Interface(keyvals[i].(string), keyvals[i+1])
	}
	e.Send()

	return nil
}

func (l *KratosLogger) Sync() error {
	return nil
}
