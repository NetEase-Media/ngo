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

package murmur3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash32(t *testing.T) {
	murmurHash := NewMurmurHash(0)
	hash := murmurHash.HashBytes([]byte("hello, world"))
	assert.Equal(t, 345750399, hash)
	hash = murmurHash.HashBytes([]byte("hello, world1"))
	assert.Equal(t, -580753116, hash)
	hash = murmurHash.HashInt32(100)
	assert.Equal(t, 616682048, hash)
	hash = murmurHash.HashInt64(100)
	assert.Equal(t, -970256272, hash)
	hash = murmurHash.HashInt32(-100)
	assert.Equal(t, 359108696, hash)
	hash = murmurHash.HashInt64(-100)
	assert.Equal(t, 16143486, hash)
}
