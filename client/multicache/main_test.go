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

package multicache

import (
	"fmt"
	"os"
	"testing"

	"github.com/NetEase-Media/ngo/client/redis"
)

func TestMain(m *testing.M) {
	fmt.Printf("multicache start.")
	// cache := gcache.New(100000).Simple().Build()
	// InitLocal(cache)

	redisOpts := []redis.Options{
		{
			Name:     "client1",
			Addr:     []string{"127.0.0.1:6379"},
			Password: "rntestncr",
			ConnType: "client",
		},
	}
	redis.Init(redisOpts)

	// redisClient := redis.GetClient("client1")
	// InitRedis(redisClient)

	opts := []Options{
		{
			Type:     "local",
			Priority: 0,
			Capacity: 100000,
		},
		{
			Type:         "redis",
			Priority:     1,
			DefaultRedis: "client1",
		},
	}
	Init(opts)
	m.Run()
	fmt.Printf("multicache end.")
	os.Exit(0)
}
