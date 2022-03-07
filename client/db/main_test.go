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
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"

	"github.com/NetEase-Media/ngo/adapter/util"
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
		Logger: New(dblogger.Config{SlowThreshold: 200 * time.Millisecond}),
	}
	testGormDB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      testMockDB,
	}), gormConfig)
	// }), &gorm.Config{PrepareStmt: true})
	util.CheckError(err)
	testClient = &Client{
		DB: testGormDB,
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
		DB: gormDB,
	}
}
