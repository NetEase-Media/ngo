package log

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	formatTXT   = "txt"
	formatJSON  = "json"
	formatBlank = "blank"

	timeFormat        = "2006-01-02 15:04:05.999"
	DefaultLoggerName = "default"
)

func New(opt *Options) (Logger, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(opt.Level)); err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var encoder zapcore.Encoder

	switch opt.Format {
	case formatJSON:
		encoder = NewJSONEncoder(opt, encoderCfg)
	case formatTXT:
		encoder = NewTextEncoder(opt, encoderCfg)
	case formatBlank:
		encoderCfg = zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "",
			NameKey:        "",
			CallerKey:      "",
			FunctionKey:    "",
			MessageKey:     "msg",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		encoder = NewTextEncoder(opt, encoderCfg)
	default:
		encoder = NewTextEncoder(opt, encoderCfg)
	}

	var core zapcore.Core
	if opt.NoFile {
		core = zapcore.NewCore(encoder, os.Stdout, lv)
	} else {
		rlog, err := newRotateLog(opt, opt.Path, "log")
		if err != nil {
			return nil, err
		}
		relog, err := newRotateLog(opt, opt.ErrorPath, "error.log")
		if err != nil {
			return nil, err
		}

		elv := zapcore.ErrorLevel
		if errlogLevel, err := zapcore.ParseLevel(opt.ErrlogLevel); err == nil {
			elv = errlogLevel
		}

		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(rlog), lv),
			zapcore.NewCore(encoder, zapcore.AddSync(relog), elv),
		)
	}

	zapOptions := make([]zap.Option, 0)
	if opt.WritableCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(opt.Skip))
	}
	if opt.WritableStack {
		zapOptions = append(zapOptions, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	zapOptions = append(zapOptions, zap.Hooks(metricsHook()))

	zl := zap.New(
		core,
		zapOptions...,
	)

	return &NgoLogger{
		lv:  &lv,
		zl:  zl,
		zsl: zl.Sugar(),
		opt: opt,
	}, nil
}

type NgoLogger struct {
	lv  *zap.AtomicLevel
	zl  *zap.Logger
	zsl *zap.SugaredLogger
	opt *Options
}

func (l *NgoLogger) WithOptions(opts ...zap.Option) *NgoLogger {
	zl := l.zl.WithOptions(opts...)
	return &NgoLogger{
		lv:  l.lv,
		zl:  zl,
		zsl: zl.Sugar(),
		opt: l.opt,
	}
}

func (l *NgoLogger) WithField(key string, value interface{}) Logger {
	zsl := l.zsl.With(key, value)
	return &NgoLogger{
		lv:  l.lv,
		zl:  zsl.Desugar(),
		zsl: zsl,
		opt: l.opt,
	}
}

func (l *NgoLogger) WithFields(key1 string, value1 interface{}, key2 string, value2 interface{}, kvs ...interface{}) Logger {
	data := make([]interface{}, 0, 4+len(kvs))
	data = append(data, key1, value1, key2, value2)
	if len(kvs) > 0 {
		data = append(data, kvs...)
	}
	zsl := l.zsl.With(data...)
	return &NgoLogger{
		lv:  l.lv,
		zl:  zsl.Desugar(),
		zsl: zsl,
		opt: l.opt,
	}
}

func (l *NgoLogger) Infof(format string, args ...interface{}) {
	l.zsl.Infof(format, args...)
}

func (l *NgoLogger) Warnf(format string, args ...interface{}) {
	l.zsl.Warnf(format, args...)
}

func (l *NgoLogger) Errorf(format string, args ...interface{}) {
	l.zsl.Errorf(format, args...)
}

func (l *NgoLogger) Debugf(format string, args ...interface{}) {
	l.zsl.Debugf(format, args...)
}

func (l *NgoLogger) Tracef(format string, args ...interface{}) {
}

func (l *NgoLogger) Fatalf(format string, args ...interface{}) {
	l.zsl.Fatalf(format, args...)
}

func (l *NgoLogger) Panicf(format string, args ...interface{}) {
	l.zsl.Panicf(format, args...)
}

func (l *NgoLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zsl.Infow(msg, keysAndValues...)
}

func (l *NgoLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zsl.Warnw(msg, keysAndValues...)
}

func (l *NgoLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zsl.Errorw(msg, keysAndValues...)
}

func (l *NgoLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zsl.Debugw(msg, keysAndValues...)
}

func (l *NgoLogger) Tracew(msg string, keysAndValues ...interface{}) {
}

func (l *NgoLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zsl.Fatalw(msg, keysAndValues...)
}

func (l *NgoLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zsl.Panicw(msg, keysAndValues...)
}

func (l *NgoLogger) Infol(msg string, fields ...Field) {
	l.zl.Info(msg, fields...)
}

func (l *NgoLogger) Warnl(msg string, fields ...Field) {
	l.zl.Warn(msg, fields...)
}

func (l *NgoLogger) Errorl(msg string, fields ...Field) {
	l.zl.Error(msg, fields...)
}

func (l *NgoLogger) Debugl(msg string, fields ...Field) {
	l.zl.Debug(msg, fields...)
}

func (l *NgoLogger) Tracel(msg string, fields ...Field) {
}

func (l *NgoLogger) Fatall(msg string, fields ...Field) {
	l.zl.Fatal(msg, fields...)
}

func (l *NgoLogger) Panicl(msg string, fields ...Field) {
	l.zl.Panic(msg, fields...)
}

func (l *NgoLogger) Info(args ...interface{}) {
	l.zsl.Info(args...)
}

func (l *NgoLogger) Warn(args ...interface{}) {
	l.zsl.Warn(args...)
}

func (l *NgoLogger) Error(args ...interface{}) {
	l.zsl.Error(args...)
}

func (l *NgoLogger) Debug(args ...interface{}) {
	l.zsl.Debug(args...)
}

func (l *NgoLogger) Trace(args ...interface{}) {
}

func (l *NgoLogger) Fatal(args ...interface{}) {
	l.zsl.Fatal(args...)
}

func (l *NgoLogger) Panic(args ...interface{}) {
	l.zsl.Panic(args...)
}

func (l *NgoLogger) GetLevel() Level {
	return l.lv.Level()
}

func (l *NgoLogger) SetLevel(lv Level) {
	l.lv.SetLevel(lv)
}

func (l *NgoLogger) Sync() {
	l.zl.Sync()
}

func metricsHook() func(zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		return nil
	}
}

func newRotateLog(opt *Options, p, suffix string) (io.Writer, error) {
	dir, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}

	linkName := path.Join(dir, opt.FileName+"."+suffix)

	return &lumberjack.Logger{
		Filename:   linkName,
		MaxSize:    opt.MaxSize, // megabytes
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge, // days
		LocalTime:  true,
		Compress:   opt.Compress, // disabled by default
	}, nil
}
