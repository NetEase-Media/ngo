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
	f := filterVersion("github.com/NetEase-Media/ngo@v0.1.52-0.20210323081228-30f33c1c0ef7/pkg/log/logger.go:119")
	assert.Equal(t, f, "github.com/NetEase-Media/ngo/pkg/log/logger.go:119")
}
