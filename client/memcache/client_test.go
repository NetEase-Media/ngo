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

package memcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var opts []Options

	o := Options{
		Name: "m1",
		Addr: []string{"127.0.0.1:11312"},
	}

	opts = append(opts, o)

	if memcacheClients == nil {
		err := Init(opts)
		assert.Equal(t, nil, err, "encounter error.")
	}

	c := GetClient("m1")
	assert.NotEqual(t, nil, c, "Init Client Failed")
}
