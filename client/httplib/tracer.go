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

import "github.com/valyala/fasthttp"

type FasthttpCarrier struct {
	header *fasthttp.RequestHeader
}

func NewFasthttpCarrier(header *fasthttp.RequestHeader) *FasthttpCarrier {
	return &FasthttpCarrier{header: header}
}

// TextMapWriter
func (carrier *FasthttpCarrier) Set(key, val string) {
	carrier.header.Set(key, val)
}

// TextMapReader
func (carrier *FasthttpCarrier) ForeachKey(handler func(key, val string) error) (err error) {
	carrier.header.VisitAll(func(k, v []byte) {
		if err = handler(string(k), string(v)); err != nil {
			return
		}
	})
	return err
}
