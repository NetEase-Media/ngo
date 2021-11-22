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

package redis

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	opts := []Options{
		{
			Name:     "client1",
			Addr:     []string{"127.0.0.1:2379"},
			ConnType: "client",
		},
		{
			Name:     "cluster1",
			Addr:     []string{"20.1.1.1", "20.1.1.2", "20.1.1.3"},
			ConnType: "cluster",
		},
	}

	err := Init(opts)
	assert.Nil(t, err)
	client := GetClient("client1").(*redisContainer)
	assert.Equal(t, []string{"127.0.0.1:2379"}, client.opt.Addr)

	cluster := GetClient("cluster1").(*redisContainer)
	assert.Equal(t, []string{"20.1.1.1", "20.1.1.2", "20.1.1.3"}, cluster.opt.Addr)

}

func generateKey() string {
	t := time.Now().Unix()
	return "test_key:" + strconv.FormatInt(t, 10) + "_" + strconv.FormatInt(int64(rand.Int()), 10)
}
