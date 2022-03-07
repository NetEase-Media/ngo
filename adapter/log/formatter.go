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
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02 15:04:05.999"
)

type Formatter struct {
	Opt *Options
}

func (formatter *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	content, err := parseEntry(entry, formatter.Opt)
	if err != nil {
		return nil, err
	}

	if isPackageLevelLessThanEntryLevel(formatter.Opt, content.Package, entry.Level) {
		return []byte{}, nil
	}

	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%s [%s] [%s] [%s]", content.Time, content.Level, content.File, content.Function)
	fmt.Fprintf(&buffer, " %s", content.Message)
	if len(entry.Data) > 0 && entry.Data[DataKey] != nil {
		fmt.Fprintf(&buffer, " %s", slice2String(entry.Data[DataKey].([]interface{})))
	}
	fmt.Fprint(&buffer, "\n")
	if content.Stack != nil {
		fmt.Fprintf(&buffer, "%s\n", content.Stack)
	}

	return buffer.Bytes(), nil
}

type JsonFormatter struct {
	Opt *Options
}

func (formatter *JsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	content, err := parseEntry(entry, formatter.Opt)
	if err != nil {
		return nil, err
	}

	if isPackageLevelLessThanEntryLevel(formatter.Opt, content.Package, entry.Level) {
		return []byte{}, nil
	}

	var msg = content.Message
	if content.Stack != nil {
		msg = fmt.Sprintf("%s\n%s", msg, content.Stack)
	}

	m := map[string]interface{}{
		"time":  content.Time,
		"level": content.Level,
		"file":  content.File,
		"func":  content.Function,
		"msg":   msg,
	}

	for k, v := range entry.Data {
		m[k] = v
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	b = append(b, '\n')
	return b, nil
}

type BlankFormatter struct {
	Opt *Options
}

func (formatter *BlankFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	content, err := parseEntry(entry, formatter.Opt)
	if err != nil {
		return nil, err
	}

	if isPackageLevelLessThanEntryLevel(formatter.Opt, content.Package, entry.Level) {
		return []byte{}, nil
	}
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%s\n", entry.Message)
	return buffer.Bytes(), nil
}

// 特定包 如果设置级别小于打印级别，那么输出空
func isPackageLevelLessThanEntryLevel(opt *Options, pkg string, entryLevel logrus.Level) bool {
	packageLevel, ok := opt.packageLogLevel[pkg]
	return ok && packageLevel < entryLevel
}
