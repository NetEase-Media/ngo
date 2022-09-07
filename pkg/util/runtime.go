package util

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// 需要忽略的几个包名
const (
	utilPkgName = "github.com/NetEase-Media/ngo/pkg/util"
)

const (
	maxCallerDepth = 20
)

var ignorePkgMap map[string]struct{}

func init() {
	ignorePkgMap = make(map[string]struct{})
	ignorePkgMap[utilPkgName] = struct{}{}
}

// CodeFrame 包含当前调用代码的属性
type CodeFrame struct {
	File     string
	Line     int
	Function string
	Package  string
}

// GetCodeFrame 或者当前调用代码的属性，用来做输出日志等调试功能
func GetCodeFrame(minDepth int, ignoreKeywords ...string) (*CodeFrame, error) {
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minDepth+2, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	var pkg, function string
	var frame *runtime.Frame
out:
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg, function = getPackageName(f.Function)

		if _, ok := ignorePkgMap[pkg]; ok {
			continue
		}

		for _, kw := range ignoreKeywords {
			if !strings.HasSuffix(f.File, "_test.go") && strings.Contains(filterVersion(f.File), kw) {
				continue out
			}
		}

		frame = &f
		break
	}

	if frame == nil {
		return nil, errors.New("can not get caller")
	}

	return &CodeFrame{
		File:     trimFile(frame.File),
		Line:     frame.Line,
		Function: function,
		Package:  pkg,
	}, nil
}

func Stack(skip int, ignoreKeywords ...string) ([]byte, error) {
	buf := new(bytes.Buffer)
out:
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		for _, kw := range ignoreKeywords {
			if !strings.HasSuffix(file, "_test.go") && strings.Contains(filterVersion(file), kw) {
				continue out
			}
		}

		if _, err := fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc); err != nil {
			break
		}
	}
	return buf.Bytes(), nil
}

func getPackageName(f string) (pkg, funtion string) {
	dotID := len(f)
	for {
		lastPeriod := strings.LastIndex(f[:dotID], ".")
		lastSlash := strings.LastIndex(f[:dotID], "/")
		if lastPeriod > lastSlash {
			dotID = lastPeriod
		} else {
			break
		}
	}

	return f[:dotID], f[dotID+1:]
}

// trimFile 将文件路径截取为最后两级
func trimFile(file string) string {
	id := len(file)
	for i := 0; i < 2; i++ {
		id = strings.LastIndex(file[:id], "/")
		if id == -1 {
			return file
		}
	}

	return file[id+1:]
}

// filterVersion 过滤版本信息
func filterVersion(file string) string {
	id1 := strings.LastIndex(file, "@v")
	if id1 == -1 {
		return file
	}
	id2 := strings.Index(file[id1:], "/")
	return file[:id1] + file[id1:][id2:]
}
