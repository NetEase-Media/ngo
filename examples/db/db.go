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
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"

	"github.com/NetEase-Media/ngo/client/db"
	"github.com/NetEase-Media/ngo/server"
)

// go run . -c ./app.yaml
func main() {
	var client *db.Client

	ctx := context.Background()
	s := server.Init()
	s.PreStart = func() error {
		client = db.GetClient("newsclient")
		var t test
		client.Raw("select blackword from blacklist t where bid = ?", 1).Find(&t)
		log.Info(t.Blacklist)

		client.Exec(context.Background(), "update blacklist set blackword = ? where bid = ?", "mmm", 21)
		return nil
	}
	go s.Start()
	defer s.Stop(ctx)

	time.Sleep(5 * time.Second)
}

type test struct {
	Blacklist string `json:"blacklist"`
}
