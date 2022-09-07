package log

import (
	"errors"

	"github.com/NetEase-Media/ngo/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options 是日志配置选项
type Options struct {
	Name            string // 代表logger
	NoFile          bool   // 是否为开发模式，如果是true，只显示到标准输出，同旧的 NoFile
	Format          string
	WritableStack   bool // 是否需要打印error及以上级别的堆栈信息
	Skip            int
	WritableCaller  bool // 是否需要打印行号函数信息
	Level           string
	Path            string
	FileName        string
	PackageLevel    map[string]string // 包级别日志等级设置
	ErrlogLevel     string            // 错误日志级别，默认error
	ErrorPath       string
	MaxAge          int  // 保留旧文件的最大天数，默认7天
	MaxBackups      int  // 保留旧文件的最大个数，默认7个
	MaxSize         int  // 在进行切割之前，日志文件的最大大小（以MB为单位）默认1024
	Compress        bool // 是否压缩/归档旧文件
	packageLogLevel map[string]Level
}

func NewDefaultOptions() *Options {
	return &Options{
		Name:           DefaultLoggerName,
		NoFile:         true,
		Format:         formatTXT,
		WritableCaller: true,
		Skip:           2,
		WritableStack:  true,
		Level:          zap.InfoLevel.String(),
		Path:           "./logs",
		FileName:       env.GetAppName(),
		PackageLevel:   make(map[string]string),
		ErrlogLevel:    zap.ErrorLevel.String(),
		ErrorPath:      "./logs",
		MaxAge:         7,
		MaxBackups:     7,
		MaxSize:        1024,
		Compress:       false,
	}
}

func checkOptions(opt *Options) error {
	if opt.Name == "" {
		return errors.New("log name can not be nil")
	}
	if len(opt.PackageLevel) > 0 && opt.WritableCaller {
		for packageStr, packageLevelStr := range opt.PackageLevel {
			packageLevel, err := zapcore.ParseLevel(packageLevelStr)
			if err == nil {
				if opt.packageLogLevel == nil {
					opt.packageLogLevel = make(map[string]zapcore.Level)
				}
				opt.packageLogLevel[packageStr] = packageLevel
			}
		}
	}
	return nil
}
