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

package cron

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"

	"github.com/stretchr/testify/assert"
)

type DummyJob struct{}

func (d DummyJob) Run() {
	panic("YOLO")
}

func TestNew(t *testing.T) {
	c := NewCron()
	assert.Equal(t, time.Local, c.Location())
	assert.Equal(t, 0, len(c.Entries()))

	c = NewCron(WithLocation(time.UTC))
	assert.Equal(t, time.UTC, c.Location())
}

func TestAddFunc(t *testing.T) {
	c := NewCron()
	entryId, _ := c.AddFunc("1 * * * *", func() { fmt.Println("print at 1") })
	assert.Equal(t, 1, len(c.Entries()))
	assert.NotNil(t, c.Entry(entryId))
}

func TestAddJob(t *testing.T) {
	var job DummyJob
	c := NewCron()
	entryId, _ := c.AddJob("1 * * * *", job)
	assert.Equal(t, 1, len(c.Entries()))
	assert.NotNil(t, c.Entry(entryId))
}

func TestRemove(t *testing.T) {
	var job DummyJob
	c := NewCron()
	entryId, _ := c.AddJob("1 * * * *", job)
	assert.Equal(t, 1, len(c.Entries()))
	assert.NotNil(t, c.Entry(entryId))
	c.Remove(entryId)
	assert.Equal(t, 0, len(c.Entries()))
}

func TestNgoLoggerWrapper(t *testing.T) {
	alogger := &NgoLoggerWrapper{logger: log.Logger()}
	alogger.Info("test", 1, 2)
	alogger.Error(errors.New("test"), "test error", 2, 3)
}
