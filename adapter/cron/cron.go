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
	"context"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"
	cron3 "github.com/robfig/cron/v3"
)

type Cron struct {
	cron *cron3.Cron
}

type Job interface {
	Run()
}

type JobWrapper func(Job) Job

type Option func(*Cron)

type EntryID cron3.EntryID

type Entry struct {
	ID   EntryID
	Next time.Time
	Prev time.Time
}

// cron 开始部分

func NewCron(opts ...Option) *Cron {
	c := &Cron{}
	c.cron = cron3.New()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Cron) Location() *time.Location {
	return c.cron.Location()
}

func (c *Cron) Start() {
	c.cron.Start()
}

func (c *Cron) Stop() context.Context {
	return c.cron.Stop()
}

func (c *Cron) AddFunc(spec string, cmd func()) (EntryID, error) {
	entryId, err := c.cron.AddFunc(spec, cmd)
	return EntryID(entryId), err
}

func (c *Cron) AddJob(spec string, cmd Job) (EntryID, error) {
	entryId, err := c.cron.AddJob(spec, cmd)
	return EntryID(entryId), err
}

func (c *Cron) Entries() []*Entry {
	entries := c.cron.Entries()
	result := make([]*Entry, 0)
	for _, e := range entries {
		e1 := &Entry{ID: EntryID(e.ID), Next: e.Next, Prev: e.Prev}
		result = append(result, e1)
	}
	return result
}

func (c *Cron) Entry(id EntryID) *Entry {
	e := c.cron.Entry(cron3.EntryID(id))
	return &Entry{ID: EntryID(e.ID), Next: e.Next, Prev: e.Prev}
}

func (c *Cron) Remove(id EntryID) {
	c.cron.Remove(cron3.EntryID(id))
}

// Cron 结束部分

// Option 开始部分
func WithLocation(loc *time.Location) Option {
	return func(c *Cron) {
		cron3.WithLocation(loc)(c.cron)
	}
}

//秒级cron支持
func WithSeconds() Option {
	return func(c *Cron) {
		cron3.WithSeconds()(c.cron)
	}
}

func WithChain(wrappers ...JobWrapper) Option {
	return func(c *Cron) {
		cron3.WithChain()(c.cron)
	}
}

// 传入logger
func WithLogger(logger *log.NgoLogger) Option {
	return func(c *Cron) {
		tlogger := &NgoLoggerWrapper{logger: logger}
		cron3.WithLogger(tlogger)(c.cron)
	}
}

// Option 结束部分

// logger 开始部分
type NgoLoggerWrapper struct {
	logger *log.NgoLogger
}

func (nlw *NgoLoggerWrapper) Info(msg string, keysAndValues ...interface{}) {
	nlw.logger.Infof("len:%d, msg:%s, values:%v", len(keysAndValues), msg, keysAndValues)
}

func (nlw *NgoLoggerWrapper) Error(err error, msg string, keysAndValues ...interface{}) {
	nlw.logger.Errorf("len:%d, msg:%s, error:%v, values:%v", len(keysAndValues), msg, err, keysAndValues)
}

// logger 结束部分
