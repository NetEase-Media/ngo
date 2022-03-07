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
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"time"
)

// A Pool maintains a pool of Redis connections.
type Pool interface {
	Get(ctx context.Context) (Conn, error)
}

type Conn interface {
	Get(name string) (string, error)
	Set(name string, value string) (bool, error)
	SetNX(name string, value string, expiry time.Duration) (bool, error)
	Eval(script *Script, keysAndArgs ...interface{}) (interface{}, error)
	PTTL(name string) (time.Duration, error)
	Close() error
}

type Script struct {
	KeyCount int
	Src      string
	Hash     string
}

func NewScript(keyCount int, src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{keyCount, src, hex.EncodeToString(h.Sum(nil))}
}
