package kafka

import (
	"github.com/NetEase-Media/ngo/pkg/log"
)

func NewLogger() *logger {
	return &logger{
		logger: log.DefaultLogger(),
	}
}

type logger struct {
	logger log.Logger
}

func (l *logger) Print(v ...interface{}) {
	l.logger.Info(v...)
}
func (l *logger) Printf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}
func (l *logger) Println(v ...interface{}) {
	l.logger.Info(v...)
}
