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

import "fmt"

var (
	_ error = NoSuchDBError{}
)

// NoSuchDBError 错误表示找不到对应name的db client
type NoSuchDBError struct {
	DBName string
}

func (err NoSuchDBError) Error() string {
	return fmt.Sprintf("can't find db named %s", err.DBName)
}
