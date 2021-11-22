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

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NetEase-Media/ngo/adapter/protocol"
	"github.com/NetEase-Media/ngo/client/db"
	"github.com/NetEase-Media/ngo/server"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type user struct {
	Id   int
	Name string
	age  int
}

func main() {
	tempdb, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("sql mock new error: %s\n", err)
		return
	}
	defer tempdb.Close()
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      tempdb,
	}), &gorm.Config{})
	if err != nil {
		fmt.Printf("gorm open error: %s\n", err)
		return
	}
	client := db.Client{DB: gdb}

	s := server.Init()
	s.AddRoute(server.GET, "/hello", func(ctx *gin.Context) {
		ctx.JSON(protocol.JsonBody("hello"))
	})
	s.AddRoute(server.GET, "/db", func(ctx *gin.Context) {
		fmt.Printf("hello, db\n")
		rows := sqlmock.NewRows([]string{"id", "name", "age"})
		for i := 0; i < 5000; i++ {
			rows = rows.AddRow(strconv.Itoa(i+1), "aaa", i+10)
		}

		mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)
		var users []user
		client.Raw("SELECT * FROM `users`").Find(&users)
		ctx.JSON(protocol.JsonBody(len(users)))
	})
	s.AddRoute(server.GET, "/db1", func(ctx *gin.Context) {
		rows := sqlmock.NewRows([]string{"id", "name", "age"})
		rows = rows.AddRow(strconv.Itoa(1), "aaa", 10)

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id`=1").WillReturnRows(rows)
		var u user
		u.Id = 1
		client.Find(context.Background(), &u)
		ctx.JSON(protocol.JsonBody(u))
	})
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("test response"))
	}))
	defer testServer.Close()

	s.Start()
}
