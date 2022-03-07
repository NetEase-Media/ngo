// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	DataKey = "data"

	formatTXT         = "txt"
	formatJSON        = "json"
	formatBlank       = "blank"
	defaultLoggerName = "default"
)

var (
	loggers map[string]*NgoLogger
)

type Fields = logrus.Fields

type NgoLogger struct {
	log logrus.Ext1FieldLogger
}

func Logger() *NgoLogger {
	return logger
}

func (l *NgoLogger) WithField(key string, value interface{}) *NgoLogger {
	var data = []interface{}{key, value}
	return &NgoLogger{
		log: l.log.WithField(DataKey, data),
	}
}

func (l *NgoLogger) WithFields(key1 string, value1 interface{}, key2 string, value2 interface{}, kvs ...interface{}) *NgoLogger {
	var data = make([]interface{}, 0, 4+len(kvs))
	data = append(data, key1, value1, key2, value2)
	if len(kvs) > 0 {
		data = append(data, kvs...)
	}
	return &NgoLogger{
		log: l.log.WithField(DataKey, data),
	}
}

func (l *NgoLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *NgoLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *NgoLogger) Errorf(format string, args ...interface{}) {
	l.withError(args...).log.Errorf(format, args...)
}

func (l *NgoLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *NgoLogger) Tracef(format string, args ...interface{}) {
	l.log.Tracef(format, args...)
}

func (l *NgoLogger) Fatalf(format string, args ...interface{}) {
	l.withError(args...).log.Fatalf(format, args...)
}

func (l *NgoLogger) Panicf(format string, args ...interface{}) {
	l.withError(args...).log.Panicf(format, args...)
}

func (l *NgoLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *NgoLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *NgoLogger) Error(args ...interface{}) {
	l.withError(args...).log.Error(args...)
}

func (l *NgoLogger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *NgoLogger) Trace(args ...interface{}) {
	l.log.Trace(args...)
}

func (l *NgoLogger) Fatal(args ...interface{}) {
	l.withError(args...).log.Fatal(args...)
}

func (l *NgoLogger) Panic(args ...interface{}) {
	l.withError(args...).log.Panic(args...)
}

func (l *NgoLogger) withError(args ...interface{}) *NgoLogger {
	if len(args) > 0 {
		if e, ok := args[len(args)-1].(error); ok {
			return &NgoLogger{
				log: l.log.WithError(e),
			}
		}
	}
	return l
}

func (l *NgoLogger) Level() logrus.Level {
	lo := l.log.(*logrus.Logger)
	return lo.Level
}

var (
	logger *NgoLogger
	// options         *Options
)

func WithField(key string, value interface{}) *NgoLogger {
	return logger.WithField(key, value)
}

func WithFields(key1 string, value1 interface{}, key2 string, value2 interface{}, kvs ...interface{}) *NgoLogger {
	return logger.WithFields(key1, value1, key2, value2, kvs...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)

}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Options 是日志配置选项
type Options struct {
	Name            string // 代表logger
	Level           string
	Path            string
	ErrorPath       string
	WritableStack   bool // 是否需要打印error及以上级别的堆栈信息
	FileName        string
	Format          string
	NoFile          bool   // 如果是true，只显示到标准输出
	FilePathPattern string // 定义文件路径名称格式
	// 包级别日志等级设置
	PackageLevel    map[string]string
	packageLogLevel map[string]logrus.Level

	// 默认7天
	MaxAge time.Duration

	// 默认1天
	RotationTime time.Duration

	// 单位MB，默认1024
	RotationSize int64
}

func NewDefaultOptions() *Options {
	return &Options{
		Path:            "",
		Level:           logrus.InfoLevel.String(),
		ErrorPath:       "",
		WritableStack:   false,
		Format:          formatTXT,
		MaxAge:          time.Hour * 24 * 7,
		RotationTime:    time.Hour * 24,
		RotationSize:    1024,
		NoFile:          true,
		PackageLevel:    make(map[string]string),
		packageLogLevel: make(map[string]logrus.Level),
	}
}

func Init(options []Options, appName string) error {
	if len(options) == 0 {
		return nil
	}
	loggers = make(map[string]*NgoLogger)
	for i := range options {
		option := &options[i]
		if len(option.Name) == 0 {
			option.Name = defaultLoggerName
		}

		if loggers[option.Name] != nil {
			return fmt.Errorf("duplicated logger config %s", option.Name)
		}

		if option.Name == defaultLoggerName && len(option.FileName) == 0 {
			option.FileName = appName
		}

		ngoLogger, err := InitLogger(option)
		if err == nil {
			loggers[option.Name] = ngoLogger
		}
	}
	if loggers[defaultLoggerName] != nil {
		logger = loggers[defaultLoggerName]
	} else {
		return fmt.Errorf("no default logger config")
	}
	return nil
}

func InitLogger(opt *Options) (*NgoLogger, error) {
	// options = opt
	level, err := logrus.ParseLevel(opt.Level)
	if err != nil {
		return nil, err
	}

	logger = &NgoLogger{log: logrus.New()}
	l := logger.log.(*logrus.Logger)

	err = setOutput(opt, l)
	if err != nil {
		return nil, err
	}

	l.SetLevel(level)

	for packageStr, packageLevelStr := range opt.PackageLevel {
		packageLevel, err := logrus.ParseLevel(packageLevelStr)
		if err == nil {
			if opt.packageLogLevel == nil {
				opt.packageLogLevel = make(map[string]logrus.Level)
			}
			opt.packageLogLevel[packageStr] = packageLevel
		}
	}

	switch opt.Format {
	case formatJSON:
		l.SetFormatter(&JsonFormatter{Opt: opt})
	case formatTXT:
		l.SetFormatter(&Formatter{Opt: opt})
	case formatBlank:
		l.SetFormatter(&BlankFormatter{Opt: opt})
	default:
		l.SetFormatter(&Formatter{Opt: opt})
	}

	return logger, nil
}

func setOutput(opt *Options, l *logrus.Logger) error {
	if opt.NoFile {
		l.SetOutput(os.Stdout)
	} else {
		// 全部日志输出
		rlAll, err := newRotateLog(opt, opt.Path, "log")
		if err != nil {
			return err
		}
		l.SetOutput(rlAll)

		// 错误日志输出
		rlError, err := newRotateLog(opt, opt.ErrorPath, "error.log")
		l.AddHook(&errorHook{writer: rlError})
	}

	// 错误日志上报哨兵
	l.AddHook(&metricsHook{Opt: opt})
	return nil
}

func newRotateLog(opt *Options, p, suffix string) (io.Writer, error) {
	dir, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}
	//有需求要自定义filepattern
	var pathPattern string
	if len(opt.FilePathPattern) > 3 {
		pathPattern = opt.FilePathPattern
	} else {
		pathPattern = path.Join(dir, opt.FileName+".%Y-%m-%d-%H-%M."+suffix)
	}
	linkName := path.Join(dir, opt.FileName+"."+suffix)
	return rotatelogs.New(
		pathPattern,
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithRotationTime(opt.RotationTime),
		rotatelogs.WithMaxAge(opt.MaxAge),
		rotatelogs.WithRotationSize(opt.RotationSize*1024*1024),
	)
}

// errorHook 将错误日志写入单独的日志中
type errorHook struct {
	writer io.Writer
}

func (h *errorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
}

func (h *errorHook) Fire(entry *logrus.Entry) error {
	b, err := entry.Bytes()
	if err != nil {
		return err
	}
	_, err = h.writer.Write(b)
	return err
}

// metricsHook 将错误日志写入哨兵
type metricsHook struct {
	Opt *Options
}

func (h *metricsHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
}

func (h *metricsHook) Fire(entry *logrus.Entry) error {
	return nil
}

func GetLogger(name string) *NgoLogger {
	return loggers[name]
}

// init 保证单测中使用的logger字段合法
func init() {
	options := []Options{*NewDefaultOptions()}
	options[0].WritableStack = true
	//options[0].PackageLevel["gorm.io/gorm"] = "info"
	Init(options, "appName")
}

func slice2String(s []interface{}) string {
	if len(s) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	fmt.Fprint(&buffer, "[")
	beg := true
	for i := 0; i < len(s); i += 2 {
		if !beg {
			fmt.Fprint(&buffer, ",")
		}
		if i == len(s)-1 {
			fmt.Fprintf(&buffer, "%s:%v", s[i], "")
		} else {
			fmt.Fprintf(&buffer, "%s:%v", s[i], s[i+1])
		}
		beg = false
	}
	fmt.Fprint(&buffer, "]")
	return buffer.String()
}
