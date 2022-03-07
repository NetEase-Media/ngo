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
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoN(t *testing.T) {
	var a, b, c int
	GoN(func() {
		a = 1
	}, func() {
		b = 2
	}, func() {
		c = 3
	})
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
	assert.Equal(t, 3, c)

	var i int32
	blockFunc := func(ctx context.Context, t int) {
		select {
		case <-ctx.Done():
			atomic.AddInt32(&i, 2)
		case <-time.NewTimer(time.Second * time.Duration(t)).C:
			atomic.AddInt32(&i, 1)
			return
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	GoN(func() {
		blockFunc(ctx, 5)
	}, func() {
		blockFunc(ctx, 10)
	}, func() {
		blockFunc(ctx, 15)
	}, func() {
		cancel()
	})

	assert.EqualValues(t, 6, i)
}
