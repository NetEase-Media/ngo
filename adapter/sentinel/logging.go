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

package sentinel

import (
	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/sirupsen/logrus"
)

func NewLogger() *Logger {
	return &Logger{
		logger: log.Logger(),
	}
}

type Logger struct {
	logging.Logger
	logger *log.NgoLogger
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if !l.DebugEnabled() {
		return
	}
	l.logger.Debug(msg, keysAndValues)
}

func (l *Logger) DebugEnabled() bool {
	return l.logger.Level() >= logrus.DebugLevel
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if !l.InfoEnabled() {
		return
	}
	l.logger.Info(msg, keysAndValues)
}
func (l *Logger) InfoEnabled() bool {
	return l.logger.Level() >= logrus.InfoLevel
}
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if !l.WarnEnabled() {
		return
	}
	l.logger.Warn(msg, keysAndValues)
}
func (l *Logger) WarnEnabled() bool {
	return l.logger.Level() >= logrus.WarnLevel
}
func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	if !l.ErrorEnabled() {
		return
	}
	l.logger.Error(msg, keysAndValues)
}
func (l *Logger) ErrorEnabled() bool {
	return l.logger.Level() >= logrus.ErrorLevel
}
