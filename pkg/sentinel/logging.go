package sentinel

import (
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/alibaba/sentinel-golang/logging"
	"go.uber.org/zap"
)

func NewLogger() *Logger {
	return &Logger{
		logger: log.DefaultLogger().(*log.NgoLogger).WithOptions(zap.AddCallerSkip(1)),
	}
}

type Logger struct {
	logging.Logger
	logger log.Logger
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if !l.DebugEnabled() {
		return
	}
	l.logger.Debug(msg, keysAndValues)
}

func (l *Logger) DebugEnabled() bool {
	return l.logger.GetLevel() >= log.DebugLevel
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if !l.InfoEnabled() {
		return
	}
	l.logger.Info(msg, keysAndValues)
}
func (l *Logger) InfoEnabled() bool {
	return l.logger.GetLevel() >= log.InfoLevel
}
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if !l.WarnEnabled() {
		return
	}
	l.logger.Warn(msg, keysAndValues)
}
func (l *Logger) WarnEnabled() bool {
	return l.logger.GetLevel() >= log.WarnLevel
}
func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	if !l.ErrorEnabled() {
		return
	}
	l.logger.Error(msg, keysAndValues)
}
func (l *Logger) ErrorEnabled() bool {
	return l.logger.GetLevel() >= log.ErrorLevel
}
