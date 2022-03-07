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

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMysqlClient(t *testing.T) {
	t.Skip("need refactor")

	opts := []*Options{
		{
			Name:            "test",
			Url:             "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
			MaxIdleCons:     10,
			MaxOpenCons:     10,
			ConnMaxLifetime: 1000,
			ConnMaxIdleTime: 10,
		},
	}
	err := Init(opts)
	assert.Nil(t, err)
}
