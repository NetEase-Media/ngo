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

package dlock

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"
	ngoredis "github.com/NetEase-Media/ngo/client/redis"
	"github.com/NetEase-Media/ngo/dlock/redis"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestDlock(t *testing.T) {
	ctx := context.Background()

	s1, _ := miniredis.Run()
	s2, _ := miniredis.Run()
	s3, _ := miniredis.Run()
	defer func() {
		s1.Close()
		s2.Close()
		s3.Close()
	}()
	opt1 := &ngoredis.Options{
		Name: fmt.Sprintf("test client %d", 1),
		Addr: []string{s1.Addr()},
	}
	rc1 := ngoredis.NewClient(opt1)
	opt2 := &ngoredis.Options{
		Name: fmt.Sprintf("test client %d", 2),
		Addr: []string{s2.Addr()},
	}
	rc2 := ngoredis.NewClient(opt2)
	opt3 := &ngoredis.Options{
		Name: fmt.Sprintf("test client %d", 3),
		Addr: []string{s3.Addr()},
	}
	rc3 := ngoredis.NewClient(opt3)
	defer func() {
		rc1.Close()
		rc2.Close()
		rc3.Close()
	}()

	dlock := initDlock(rc1, rc2, rc3)
	for j := 0; j < 3; j++ {
		testRun(ctx, dlock, 100, t)
	}
}

func TestRenew(t *testing.T) {
	ctx := context.Background()

	s1, _ := miniredis.Run()
	s2, _ := miniredis.Run()
	s3, _ := miniredis.Run()
	defer func() {
		s1.Close()
		s2.Close()
		s3.Close()
	}()
	opt1 := &ngoredis.Options{
		Name: fmt.Sprintf("test client %d", 1),
		Addr: []string{s1.Addr()},
	}
	rc1 := ngoredis.NewClient(opt1)
	opt2 := &ngoredis.Options{
		Name: fmt.Sprintf("test client %d", 2),
		Addr: []string{s2.Addr()},
	}
	rc2 := ngoredis.NewClient(opt2)
	opt3 := &ngoredis.Options{
		Name: fmt.Sprintf("test client %d", 3),
		Addr: []string{s3.Addr()},
	}
	rc3 := ngoredis.NewClient(opt3)
	defer func() {
		rc1.Close()
		rc2.Close()
		rc3.Close()
	}()

	dlock := initDlock(rc1, rc2, rc3)
	succ, executed, err := dlock.NewMutex("test", func() {
		log.Info("start working...")
		time.Sleep(time.Second * 2)
		log.Info("end working...")
	}).WithExpiry(time.Second).DoContext(ctx)

	assert.True(t, succ)
	assert.True(t, executed)
	assert.NoError(t, err)
}

func testRun(ctx context.Context, dlock *Dlock, k int, t *testing.T) {
	var (
		wg sync.WaitGroup
		n  int
	)
	for i := 0; i < k; i++ {
		wg.Add(1)
		go func() {
			succ, executed, err := dlock.NewMutex("test", func() {
				n++
			}).WithTries(100).WithRetryDelayFunc(func(tries int) time.Duration {
				return time.Duration(rand.Intn(1000-10)+10) * time.Millisecond
			}).DoContext(ctx)

			assert.True(t, succ)
			assert.True(t, executed)
			assert.NoError(t, err)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, k, n)
}

func initDlock(rcs ...ngoredis.Redis) *Dlock {
	pools := make([]redis.Pool, len(rcs))
	for i, rc := range rcs {
		pool := redis.NewPool(rc)
		pools[i] = pool
	}
	dlock := New(pools...)
	return dlock
}
