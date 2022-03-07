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

package httplib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestFasthttpCarrier(t *testing.T) {
	header := &fasthttp.RequestHeader{}
	carrier := NewFasthttpCarrier(header)
	carrier.Set("Name", "dahai")
	carrier.Set("Age", "17")
	keys := make(map[string]string)
	carrier.ForeachKey(func(k, v string) error {
		keys[k] = v
		return nil
	})
	assert.Equal(t, "dahai", keys["Name"])
	assert.Equal(t, "17", keys["Age"])
}
