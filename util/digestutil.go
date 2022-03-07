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
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// md5，返回入参的md5值
// 如果是空串，md5也返回空串
func Md5(text string) string {
	if len(text) == 0 {
		return ""
	}
	sum := md5.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}

// sha1，返回入参的sha1值
// 如果是空串，sha1也返回空串
func Sha1(text string) string {
	if len(text) == 0 {
		return ""
	}
	sum := sha1.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}
