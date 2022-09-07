package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutBoundIP(t *testing.T) {
	ip, _ := GetOutBoundIP()
	assert.NotEmpty(t, ip)
}
