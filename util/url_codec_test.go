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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	query := "a=100&b=100&c=a<"
	expect := "a%3D100%26b%3D100%26c%3Da%3C"
	actual := Encode(query)
	assert.Equal(t, expect, actual, "编码失败")
}

func TestDecode(t *testing.T) {
	query := "a=100&b=100&c=a%3C"
	expect := "a=100&b=100&c=a<"
	actual := Decode(query)
	assert.Equal(t, expect, actual, "解码失败")
}

func TestEncode01(t *testing.T) {
	query := ""
	expect := ""
	actual := Encode(query)
	assert.Equal(t, expect, actual, "编码kong失败")
}

func TestDecode01(t *testing.T) {
	query := ""
	expect := ""
	actual := Decode(query)
	assert.Equal(t, expect, actual, "解码kong失败")
}
