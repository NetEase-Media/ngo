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
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	opt := &Options{
		Path:          "./log",
		Level:         logrus.InfoLevel.String(),
		ErrorPath:     "",
		WritableStack: true,
		Format:        formatTXT,
		MaxAge:        time.Hour * 24 * 7,
		RotationTime:  time.Hour * 24,
		RotationSize:  1024,
		NoFile:        true,
	}

	opts := []Options{*opt}
	err := Init(opts, "default")

	assert.Nil(t, err)

	defer func() {
		err := recover()
		assert.Nil(t, err)
	}()
	Trace("Trace")
	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Error("Error")
	Error("Error", " ", errors.New("test error"))

	Tracef("%s", "Trace")
	Debugf("%s", "Debug")
	Infof("%s", "Info")
	Warnf("%s", "Warn")
	Errorf("%s", "Error")
	Errorf("%s %s", "Error", errors.New("test error"))

	WithField("k1", "v1").Info("into")
	WithFields("k1", "v1", "k2", "v2").Info("into")
	WithFields("k1", "v1", "k2", "v2", "k3").Info("into")
	WithFields("k1", "v1", "k2", "v2", "k3", "k4").Info("into")
	WithFields("k1", "v1", "k2", "v2").Errorf("error: %v", errors.New("test error"))
}

func BenchmarkLogMetrics1(b *testing.B) {
	initWithOpts(false)
	for i := 0; i < b.N; i++ {
		Error("error")
	}
}
func BenchmarkLogMetrics2(b *testing.B) {
	initWithOpts(true)
	for i := 0; i < b.N; i++ {
		Error("error")
	}
}
func BenchmarkLogMetrics3(b *testing.B) {
	initWithOpts(false)
	for i := 0; i < b.N; i++ {
		Error("error")
	}
}
func BenchmarkLogMetrics4(b *testing.B) {
	initWithOpts(true)
	for i := 0; i < b.N; i++ {
		Error("error")
	}
}

func initWithOpts(writableStack bool) {
	opt := &Options{
		Path:          "./log",
		Level:         logrus.InfoLevel.String(),
		ErrorPath:     "",
		WritableStack: writableStack,
		Format:        formatTXT,
		MaxAge:        time.Hour * 24 * 7,
		RotationTime:  time.Hour * 24,
		RotationSize:  1024,
		NoFile:        true,
	}

	opts := []Options{*opt}
	Init(opts, "")
}
