package db

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NetEase-Media/ngo/pkg/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
)

var (
	testMockDB  *sql.DB
	testSqlmock sqlmock.Sqlmock
	testGormDB  *gorm.DB
	testClient  *Client
)

type testuser struct {
	ID     string
	Name   string
	Gender string
}

func TestMain(m *testing.M) {
	setupTest()
	ret := m.Run()
	tearDownTest()
	os.Exit(ret)
}

func setupTest() {
	var err error
	testMockDB, testSqlmock, err = sqlmock.New()
	util.CheckError(err)
	gormConfig := &gorm.Config{
		Logger: NewLogger(dblogger.Config{SlowThreshold: 200 * time.Millisecond}),
	}
	testGormDB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      testMockDB,
	}), gormConfig)
	// }), &gorm.Config{PrepareStmt: true})
	util.CheckError(err)
	testClient = &Client{
		db: testGormDB,
	}
}

func tearDownTest() {
	testMockDB.Close()
}

func testNewORM(t *testing.T) (sqlmock.Sqlmock, *sql.DB, *Client) {
	testMockDB, testSqlmock, err := sqlmock.New()
	assert.Nil(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      testMockDB,
	}), &gorm.Config{})
	assert.Nil(t, err)
	return testSqlmock, testMockDB, &Client{
		db: gormDB,
	}
}
