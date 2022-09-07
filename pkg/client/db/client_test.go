package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gorm.io/gorm/schema"

	"github.com/stretchr/testify/assert"
)

type test struct {
	Id   int64
	Name string
}

func TestClient(t *testing.T) {
	patches := gomonkey.ApplyFunc(mysql.Open, func(dsn string) gorm.Dialector {
		db, _, _ := sqlmock.New()
		return mysql.New(mysql.Config{
			DSN:                       dsn,
			SkipInitializeWithVersion: true,
			Conn:                      db,
		})
	})
	defer patches.Reset()

	c, err := New(&Options{
		Name:            "test",
		Type:            "mysql",
		Url:             "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		MaxIdleCons:     10,
		MaxOpenCons:     10,
		ConnMaxLifetime: 1000,
		ConnMaxIdleTime: 10,
	})
	assert.Nil(t, err)
	ctx := context.Background()
	db := WithContext(ctx, c)
	db.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}
	var g test
	db.Create(&g)
	db.Find(&g)
}
