package clog

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"time"
)

const (
	errKey = "err"
)

var _ log.Logger = (*KratosLogger)(nil)
var devLogger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

type KratosLogger struct {
	log *zerolog.Logger
	dev bool
}

func NewKratosLogger(zeroLogger *zerolog.Logger, dev bool) *KratosLogger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
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

	fun := func(level log.Level, e *zerolog.Event) {
		switch level {
		case log.LevelDebug, log.LevelInfo, log.LevelWarn:
			for i := 0; i < len(keyvals); i += 2 {
				e = e.Interface(keyvals[i].(string), keyvals[i+1])
			}
		case log.LevelError, log.LevelFatal:
			keyMap := make(map[string]int)
			for i := 0; i < len(keyvals); i += 2 {
				key := keyvals[i].(string)
				keyMap[key] = i
			}
			if index, ok := keyMap[errKey]; ok {
				e = e.Stack().Err(keyvals[index+1].(error))
				delete(keyMap, errKey)
			}
			for key, keyIndex := range keyMap {
				e = e.Interface(key, keyvals[keyIndex+1])
			}
		}
		e.Send()
	}

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
			e = devLogger.Error()
		case log.LevelFatal:
			e = devLogger.Fatal()
		}

		fun(level, e)
	}

	switch level {
	case log.LevelDebug:
		e = l.log.Debug()
	case log.LevelInfo:
		e = l.log.Info()
	case log.LevelWarn:
		e = l.log.Warn()
	case log.LevelError:
		e = l.log.Error()
	case log.LevelFatal:
		e = l.log.Fatal()
	}

	fun(level, e)

	return nil
}

func (l *KratosLogger) Sync() error {
	return nil
}
