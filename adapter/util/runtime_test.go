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

func TestGetCodeFrame(t *testing.T) {
	_, err := GetCodeFrame(0)
	assert.Nil(t, err)
}

func TestFilterVersion(t *testing.T) {
	f := filterVersion("github.com/NetEase-Media/ngo@v0.1.52-0.20210323081228-30f33c1c0ef7/adapter/log/logger.go:119")
	assert.Equal(t, f, "github.com/NetEase-Media/ngo/adapter/log/logger.go:119")
}
