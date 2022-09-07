package log

import (
	"sync"
)

var defaultLogger = newDefaultLogger()

func newDefaultLogger() Logger {
	opt := NewDefaultOptions()
	logger, _ := New(opt)
	return logger
}

func DefaultLogger() Logger {
	return defaultLogger
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func WithField(key string, value interface{}) Logger {
	return defaultLogger.WithField(key, value)
}

func WithFields(key1 string, value1 interface{}, key2 string, value2 interface{}, kvs ...interface{}) Logger {
	return defaultLogger.WithFields(key1, value1, key2, value2, kvs...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Tracef(format string, args ...interface{}) {
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func Tracew(msg string, keysAndValues ...interface{}) {
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

func Infol(msg string, fields ...Field) {
	defaultLogger.Infol(msg, fields...)
}

func Warnl(msg string, fields ...Field) {
	defaultLogger.Warnl(msg, fields...)
}

func Errorl(msg string, fields ...Field) {
	defaultLogger.Errorl(msg, fields...)
}

func Debugl(msg string, fields ...Field) {
	defaultLogger.Debugl(msg, fields...)
}

func Tracel(msg string, fields ...Field) {
}

func Fatall(msg string, fields ...Field) {
	defaultLogger.Fatall(msg, fields...)
}

func Panicl(msg string, fields ...Field) {
	defaultLogger.Panicl(msg, fields...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Trace(args ...interface{}) {
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Sync() {
	defaultLogger.Sync()
}

func GetLevel() Level {
	return defaultLogger.GetLevel()
}

func SetLevel(lv Level) {
	defaultLogger.SetLevel(lv)
}

var (
	mu      sync.RWMutex
	loggers = make(map[string]Logger)
)

func SetLogger(name string, logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	loggers[name] = logger
}

func GetLogger(name string) Logger {
	mu.RLock()
	defer mu.RUnlock()
	return loggers[name]
}

func GetLoggers() map[string]Logger {
	mu.RLock()
	defer mu.RUnlock()
	return loggers
}
