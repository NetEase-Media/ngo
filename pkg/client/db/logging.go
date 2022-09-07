package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/NetEase-Media/ngo/pkg/log"
	dblogger "gorm.io/gorm/logger"
)

var (
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s [%.3fms] [rows:%v] %s"
)

func NewLogger(config dblogger.Config) *logger {
	return &logger{
		Config: config,
		logger: log.DefaultLogger().(*log.NgoLogger).WithOptions(zap.AddCallerSkip(1)),
	}
}

type logger struct {
	dblogger.Interface
	dblogger.Config
	logger log.Logger
}

// LogMode log mode
func (l *logger) LogMode(level dblogger.LogLevel) dblogger.Interface {
	return l
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	l.logger.Infof(msg, data...)
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	l.logger.Warnf(msg, data...)
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	l.logger.Errorf(msg, data...)
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logger.GetLevel() <= log.ErrorLevel:
		sql, rows := fc()
		if rows == -1 {
			l.logger.Errorf(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logger.Errorf(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logger.GetLevel() <= log.WarnLevel:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.logger.Warnf(traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logger.Warnf(traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.logger.GetLevel() <= log.InfoLevel:
		sql, rows := fc()
		if rows == -1 {
			l.logger.Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logger.Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
