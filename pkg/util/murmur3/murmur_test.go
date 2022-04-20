package murmur3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash32(t *testing.T) {
	murmurHash := NewMurmurHash(0)
	hash := murmurHash.HashBytes([]byte("hello, world"))
	assert.Equal(t, 345750399, hash)
	hash = murmurHash.HashBytes([]byte("hello, world1"))
	assert.Equal(t, -580753116, hash)
	hash = murmurHash.HashInt32(100)
	assert.Equal(t, 616682048, hash)
	hash = murmurHash.HashInt64(100)
	assert.Equal(t, -970256272, hash)
	hash = murmurHash.HashInt32(-100)
	assert.Equal(t, 359108696, hash)
	hash = murmurHash.HashInt64(-100)
	assert.Equal(t, 16143486, hash)
}
