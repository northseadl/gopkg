package log

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

var _ cron.Logger = (*CronLogger)(nil)

type CronLogger struct {
	log *zerolog.Logger
}

func NewCronLogger(zeroLogger *zerolog.Logger) *CronLogger {
	return &CronLogger{zeroLogger}
}

func (l CronLogger) Info(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 || len(keysAndValues)%2 != 0 {
		l.log.Warn().Msg(fmt.Sprint("Key-values must appear in pairs: ", keysAndValues))
		return
	}
	e := l.log.Info()
	for i := 0; i < len(keysAndValues); i = i + 2 {
		e.Interface(keysAndValues[i].(string), keysAndValues[i+1])
	}
	e.Msg(msg)
}

func (l CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 || len(keysAndValues)%2 != 0 {
		l.log.Warn().Msg(fmt.Sprint("Key-values must appear in pairs: ", keysAndValues))
		return
	}
	e := l.log.Error().Err(err)
	for i := 0; i < len(keysAndValues); i = i + 2 {
		e.Interface(keysAndValues[i].(string), keysAndValues[i+1])
	}
	e.Msg(msg)
}
