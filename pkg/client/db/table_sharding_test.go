package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestTableSharding(t *testing.T) {
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
		Url:             "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		MaxIdleCons:     10,
		MaxOpenCons:     10,
		ConnMaxLifetime: 1000,
		ConnMaxIdleTime: 10,
	})
	assert.Nil(t, err)
	var g test
	tn := xxTable("xxx")
	WithContext(context.Background(), c).Table(tn).Raw("select * from "+tn+" where id = ?", 1).Find(&g)
}

func TestTableName(t *testing.T) {
	name := NewTableSharding(WithKey("YDJ0996E2IEEFZYZ"), WithName("re_user_recommend"), WithSize(128)).TableName()
	assert.Equal(t, "re_user_recommend_7", name)
}

func xxTable(key string) string {
	return NewTableSharding(WithKey(key), WithName("test"), WithSize(8)).TableName()
}
