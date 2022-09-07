package memcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	o := Options{
		Name: "m1",
		Addr: []string{"127.0.0.1:11312"},
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "encounter error.")
	assert.NotEqual(t, nil, c, "Init Client Failed")
}
