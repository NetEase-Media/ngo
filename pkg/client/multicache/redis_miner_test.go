package multicache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisSet(t *testing.T) {
	client := GetRedisMiner()
	op, err := client.Set("kk", "kiko")
	assert.Equal(t, true, op, "")
	assert.Equal(t, true, err == nil, "")

	ret, err := client.Get("kk")
	assert.Equal(t, true, err == nil, "")
	assert.Equal(t, "kiko", ret, "")

	op, err = client.Evict("kk")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}

func TestRedisSetWithTimeout(t *testing.T) {
	client := GetRedisMiner()
	op, err := client.SetWithTimeout("kko", "kiko", 5)

	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")

	ret, err := client.Get("kko")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, "kiko", ret, "")

	time.Sleep(6 * time.Second)
	ret, err = client.Get("kko")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, "", ret, "")
}

func TestRedisClear(t *testing.T) {
	client := GetRedisMiner()
	op, err := client.Clear()
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}
