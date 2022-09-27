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
