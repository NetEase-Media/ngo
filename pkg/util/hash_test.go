package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMurmurHashString(t *testing.T) {
	hash := MurmurHashString("p:i:a:newsclient:d:{E1C734C9-947C-46FF-AA75-01EF46435899}:h")
	assert.Equal(t, int64(134429908068854083), hash)

	hash = MurmurHashString("p:i:a:newsclient:d:{693A4B76-8D5A-4B3B-8F9F-6434C5478706}:h")
	assert.Equal(t, int64(-3579132390707086426), hash)
}
