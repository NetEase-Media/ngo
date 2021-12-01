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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocalSet(t *testing.T) {
	client := GetLocalMiner()
	op, err := client.Set("kk", "kiko")
	assert.Equal(t, true, op, "")
	assert.Equal(t, true, err == nil, "")

	ret, err := client.Get("kk")
	assert.Equal(t, true, err == nil, "")
	assert.Equal(t, "kiko", ret, "")

	op, err = client.Evict("kk")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}

func TestLocalSetWithTimeout(t *testing.T) {
	client := GetLocalMiner()
	op, err := client.SetWithTimeout("kko", "kiko", 10)

	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")

	ret, err := client.Get("kko")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, "kiko", ret, "")

	time.Sleep(10 * time.Second)
	ret, err = client.Get("kko")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, "", ret, "")
}

func TestLocalClear(t *testing.T) {
	client := GetLocalMiner()
	op, err := client.Clear()
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}
