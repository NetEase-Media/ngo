package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeName(t *testing.T) {
	type testTypeA struct{}

	var a testTypeA
	name := "util.testTypeA"
	assert.Equal(t, name, TypeName(a))
	assert.Equal(t, name, TypeName(&a))

	p := &a
	pp := &p
	assert.Equal(t, name, TypeName(pp))
}
