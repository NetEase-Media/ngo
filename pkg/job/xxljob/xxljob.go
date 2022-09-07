package xxljob

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/NetEase-Media/ngo/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/xxl-job/xxl-job-executor-go"
)

var gXxlExecutor *XxlExecutor

type Options struct {
	Enabled      bool
	ServerAddr   string //调度中心地址
	AccessToken  string //请求令牌
	ExecutorIp   string //本地(执行器)IP(可自行获取)
	ExecutorPort string //本地(执行器)端口
	RegistryKey  string //执行器名称
	LogDir       string //日志目录
}

type NgoTaskFunc func(cxt context.Context, param *xxl.RunReq, logger *XxlJobLogger) string

type XxlExecutor struct {
	xxl.Executor
	LogDir string
}

type XxlJobLogger struct {
	logger *logrus.Logger
}

func NewXxlJobLogger(logid int64) (*XxlJobLogger, error) {
	file, err := getLogFile(logid)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetOutput(file)
	xxlJobLogger := &XxlJobLogger{logger: logger}
	return xxlJobLogger, nil
}

func (l *XxlJobLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *XxlJobLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func NewXxlExecutor(executor xxl.Executor, logDir string) *XxlExecutor {
	return &XxlExecutor{Executor: executor, LogDir: logDir}
}

func getLogFile(logid int64) (*os.File, error) {
	logDir := GetXxlExecutor().LogDir
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(logDir, strconv.FormatInt(logid, 10)+".log")
	fileflag := os.O_RDWR | os.O_CREATE | os.O_APPEND
	return os.OpenFile(path, fileflag, 0755)
}

func Init(opt *Options, clusterName string) (err error) {
	if !opt.Enabled {
		return nil
	}
	xxlOptions := make([]xxl.Option, 0)
	if opt.ServerAddr == "" {
		return errors.New("xxljob.serverAddr is empty")
	}
	xxlOptions = append(xxlOptions, xxl.ServerAddr(opt.ServerAddr))
	if opt.ExecutorIp == "" {
		opt.ExecutorIp, err = util.GetOutBoundIP()
		if err != nil {
			return errors.New(fmt.Sprintf("xxljob get ip error,%s", err))
		}
	}
	xxlOptions = append(xxlOptions, xxl.ExecutorIp(opt.ExecutorIp))
	if opt.ExecutorPort == "" {
		opt.ExecutorPort = "19876"
	}
	xxlOptions = append(xxlOptions, xxl.ExecutorPort(opt.ExecutorPort))
	if opt.RegistryKey == "" {
		opt.RegistryKey = clusterName
	}
	xxlOptions = append(xxlOptions, xxl.RegistryKey(opt.RegistryKey))

	if opt.LogDir == "" {
		opt.LogDir = "./log/xxljob"
	}

	executor := xxl.NewExecutor(xxlOptions...)
	gXxlExecutor = NewXxlExecutor(executor, opt.LogDir)
	gXxlExecutor.LogHandler(xxlJobLogHandler)
	gXxlExecutor.Init()
	go gXxlExecutor.Run()
	return nil
}

//注册移除
func Stop() {
	if gXxlExecutor != nil && gXxlExecutor.Executor != nil {
		gXxlExecutor.Executor.Stop()
	}
}

func GetXxlExecutor() *XxlExecutor {
	return gXxlExecutor
}

//注册任务
func RegTask(taskname string, task NgoTaskFunc) {
	gXxlExecutor.RegTask(taskname, func(ctx context.Context, param *xxl.RunReq) string {
		logger, err := NewXxlJobLogger(param.LogID)
		if err != nil {
			fmt.Printf("new xxl job logger error:%v\n", err)
			return err.Error()
		}
		return task(ctx, param, logger)
	})
}

func xxlJobLogHandler(req *xxl.LogReq) *xxl.LogRes {
	file, err := getLogFile(req.LogID)
	if err != nil {
		return nil
	}
	defer file.Close()
	rowCount := 0
	s := bufio.NewScanner(file)
	var buf bytes.Buffer
	for s.Scan() {
		rowCount++
		line := s.Text()
		if rowCount >= req.FromLineNum {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
	return &xxl.LogRes{Code: 200, Msg: "OK", Content: xxl.LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   rowCount,
		LogContent:  buf.String(),
		IsEnd:       true,
	}}
}
