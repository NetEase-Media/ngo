package db

import (
	"context"
	"time"

	"github.com/NetEase-Media/ngo/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
)

const (
	dbTypeMysql = "mysql"
)

func New(opt *Options) (*Client, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}

	var cfg gorm.Config
	cfg.Logger = NewLogger(dblogger.Config{
		SlowThreshold: 200 * time.Millisecond,
	})
	var dialector gorm.Dialector
	switch opt.Type {
	case dbTypeMysql:
		dialector = mysql.Open(opt.Url)
	default:
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

	db.Use(newGormTracerPlugin())

	client := &Client{
		db:  db,
		Opt: *opt,
	}
	return client, nil
}

// Client
type Client struct {
	db  *gorm.DB
	Opt Options
}

func (client *Client) Close() {
	if db, err := client.db.DB(); err == nil {
		db.Close()
	}
}

// WithContext
func WithContext(ctx context.Context, client *Client) *gorm.DB {
	return client.db.WithContext(ctx)
}
