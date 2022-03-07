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

package db

import (
	"context"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
)

const (
	dbTypeMysql = "mysql"
)

// Options 是MysqlClient的配置数据
type Options struct {
	Name            string
	Type            string
	Url             string
	MaxIdleCons     int
	MaxOpenCons     int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewDefaultOptions() *Options {
	return &Options{
		Type:            "mysql",
		MaxIdleCons:     10,
		MaxOpenCons:     10,
		ConnMaxLifetime: time.Second * 1000,
		ConnMaxIdleTime: time.Second * 60,
	}
}

// Client
type Client struct {
	*gorm.DB

	opt Options
}

func NewClient(opt *Options) (*Client, error) {
	var cfg gorm.Config
	cfg.Logger = New(dblogger.Config{
		SlowThreshold: 200 * time.Millisecond,
	})
	var dialector gorm.Dialector
	if opt.Type == dbTypeMysql {
		dialector = mysql.Open(opt.Url)
	} else {
		dialector = mysql.Open(opt.Url)
	}

	db, err := gorm.Open(dialector, &cfg)
	if err != nil {
		log.Errorf("can not be open client. msg:%s", err.Error())
		return nil, err
	}

	myDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	myDB.SetMaxIdleConns(opt.MaxIdleCons)
	myDB.SetMaxOpenConns(opt.MaxOpenCons)
	myDB.SetConnMaxLifetime(opt.ConnMaxLifetime)
	myDB.SetConnMaxIdleTime(opt.ConnMaxIdleTime)

	db.Use(newGormMetricsPlugin())
	db.Use(newGormTracerPlugin())

	client := &Client{
		DB:  db,
		opt: *opt,
	}
	return client, nil
}

func (client *Client) Trace(context context.Context, f func() *gorm.DB) (tx *gorm.DB) {
	span, ctx := opentracing.StartSpanFromContext(context, "gorm", ext.SpanKindRPCClient)
	oldContext := client.DB.Statement.Context
	client.DB.Statement.Context = ctx
	defer func() {
		span.Finish()
		client.DB.Statement.Context = oldContext
	}()

	return f()
}

func (client *Client) Create(context context.Context, value interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Create(value)
}
func (client *Client) Save(context context.Context, value interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Save(value)
}
func (client *Client) First(context context.Context, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).First(dest, conds)
}
func (client *Client) Take(context context.Context, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Take(dest, conds...)
}
func (client *Client) Last(context context.Context, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Last(dest, conds...)
}
func (client *Client) Find(context context.Context, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Find(dest, conds...)
}
func (client *Client) FindInBatches(context context.Context, dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) (tx *gorm.DB) {
	return client.WithContext(context).FindInBatches(dest, batchSize, fc)
}
func (client *Client) FirstOrInit(context context.Context, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).FirstOrInit(dest, conds...)
}
func (client *Client) FirstOrCreate(context context.Context, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).FirstOrCreate(dest, conds...)
}
func (client *Client) Update(context context.Context, column string, value interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Update(column, value)
}
func (client *Client) Updates(context context.Context, values interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Updates(values)
}
func (client *Client) UpdateColumn(context context.Context, column string, value interface{}) (tx *gorm.DB) {
	return client.WithContext(context).UpdateColumn(column, value)
}
func (client *Client) UpdateColumns(context context.Context, values interface{}) (tx *gorm.DB) {
	return client.WithContext(context).UpdateColumns(values)
}
func (client *Client) Delete(context context.Context, value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Delete(value, conds...)
}
func (client *Client) Count(context context.Context, count *int64) (tx *gorm.DB) {
	return client.WithContext(context).Count(count)
}
func (client *Client) Scan(context context.Context, dest interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Scan(dest)
}
func (client *Client) Pluck(context context.Context, column string, dest interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Pluck(column, dest)
}
func (client *Client) Exec(context context.Context, sql string, values ...interface{}) (tx *gorm.DB) {
	return client.WithContext(context).Exec(sql, values...)
}

func setError(err error) {
	// TODO 出现异常的时候需要做的处理
}

func duration(startTime int64) {
	// TODO 上报接口使用时长
}
