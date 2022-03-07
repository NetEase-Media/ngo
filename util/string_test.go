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

func TestJoin(t *testing.T) {
	input := []string{" as", " ", "de ", " as "}
	symbol := Colon
	expect := "as: :de : as"
	actual := Join(input, symbol)
	assert.Equal(t, expect, actual, "join failed!")
}

func TestSplit(t *testing.T) {
	input := "as: :de : as"
	actual := Split(input, Colon)
	expect := []string{"as", "de", "as"}
	assert.Equal(t, expect, actual)
}

func TestJoin01(t *testing.T) {
	input := []string{" a s ", "", "de ", " as "}
	expect := "a s __de _ as"
	actual := Join(input, Underline)
	assert.Equal(t, expect, actual, "join01 failed!")
}

func TestJoin02(t *testing.T) {
	input := []string{"   a s   ", "", "de ", " as   "}
	expect := "a s   ,,de , as"
	actual := Join(input, Comma)
	assert.Equal(t, expect, actual, "join01 failed!")
}

func TestSplit03(t *testing.T) {
	input := " add , ad, ,, ads   "
	actual := Split(input, Comma)
	expect := []string{"add", "ad", "ads"}
	assert.Equal(t, expect, actual)
}

func TestSplit04(t *testing.T) {
	input := " add | ad| || ads   "
	actual := Split(input, VerticalBar)
	expect := []string{"add", "ad", "ads"}
	assert.Equal(t, expect, actual)
}

func TestSplit05(t *testing.T) {
	input := " add | ad| || ad   "
	actual := SplitNoRepeat(input, VerticalBar)
	expect := []string{"add", "ad"}
	assert.Equal(t, expect, actual)
}
