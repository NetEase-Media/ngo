package cron

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/pkg/log"

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
	alogger := &NgoLoggerWrapper{logger: log.DefaultLogger()}
	alogger.Info("test", 1, 2)
	alogger.Error(errors.New("test"), "test error", 2, 3)
}
