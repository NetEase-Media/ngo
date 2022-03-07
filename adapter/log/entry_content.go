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
	"context"
	"strconv"
	"time"

	"github.com/NetEase-Media/ngo/adapter/util"
	"github.com/sirupsen/logrus"
)

type entryContentKey string

const (
	ecKey              entryContentKey = "entryContent"
	minimumCallerDepth int             = 0
	stackSkip          int             = 0
)

var ignoreKeywords = []string{"/ngo/adapter/log", "/ngo/adapter/sentinel/logging.go", "/ngo/client/kafka/logging.go", "/ngo/client/db/logging.go", "/ngo/client/db/ddb/driver/logging.go", "github.com/sirupsen/logrus"}

func parseEntry(entry *logrus.Entry, opt *Options) (*entryContent, error) {
	if entry.Context != nil {
		return entry.Context.Value(ecKey).(*entryContent), nil
	}

	codeFrame, err := util.GetCodeFrame(minimumCallerDepth, ignoreKeywords...)
	if err != nil {
		return nil, err
	}

	var s []byte
	// error级别以上，打印堆栈
	if entry.Level <= logrus.ErrorLevel && opt.WritableStack {
		s, _ = util.Stack(stackSkip, ignoreKeywords...)
	}

	ec := &entryContent{
		Time:     time.Now().Format(timeFormat),
		Level:    entry.Level.String(),
		File:     codeFrame.File + ":" + strconv.Itoa(codeFrame.Line),
		Function: codeFrame.Function,
		Package:  codeFrame.Package,
		Data:     entry.Data,
		Message:  entry.Message,
		Stack:    s,
	}
	// 存入上下文，给后续hook使用
	if entry.Context == nil {
		entry.Context = context.WithValue(context.Background(), ecKey, ec)
	}
	return ec, nil
}

// entryContent 是entry的内容解析结构
type entryContent struct {
	Time     string
	Level    string
	File     string
	Function string
	Data     Fields
	Message  string
	Stack    []byte
	Package  string
}
