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
	"github.com/NetEase-Media/ngo/client/redis"
	"github.com/NetEase-Media/ngo/server"
)

var (
	client redis.Redis
)

// go run . -c ./app.yaml
func main() {
	s := server.Init()

	s.PreStart = func() error {
		key := "key"
		value := "value"

		client = redis.GetClient("redis01")
		_, err := client.Set(context.Background(), key, value, time.Second*5).Result()
		if err != nil {
			log.Error(err)
		}
		res, err := client.Get(context.Background(), key).Result()
		if err != nil {
			log.Error(err)
		}
		log.Info(res)
		return nil
	}
	s.Start()
}
