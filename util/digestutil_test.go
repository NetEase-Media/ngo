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

func TestMd5(t *testing.T) {
	assert.Equal(t, "", Md5(""))
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e", Md5("123456"))
	assert.Equal(t, "fc5e038d38a57032085441e7fe7010b0", Md5("helloworld"))
	assert.Equal(t, "c8b0ff27a844d2eecd81669dbaa544eb", Md5("你好中国"))
}

func TestSha1(t *testing.T) {
	assert.Equal(t, "", Md5(""))
	assert.Equal(t, "7c4a8d09ca3762af61e59520943dc26494f8941b", Sha1("123456"))
	assert.Equal(t, "6adfb183a4a2c94a2f92dab5ade762a47889a5a1", Sha1("helloworld"))
	assert.Equal(t, "e831d9bded4675af26ceb1a35f3f86261c584392", Sha1("你好中国"))
}
