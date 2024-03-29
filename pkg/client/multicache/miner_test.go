package multicache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	miner, err := GetMiner()

	assert.EqualValues(t, nil, err, "")
	if miner == nil {
		return
	}

	redisClient := GetRedisMiner()
	op, err := redisClient.SetWithTimeout("kk", "kiko", 20)
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")

	ret, err := miner.Get("kk")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, "kiko", ret, "")
}

func TestSet(t *testing.T) {
	miner, err := GetMiner()

	assert.EqualValues(t, nil, err, "")
	if miner == nil {
		return
	}
	op, err := miner.Set("kk", "kiko")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}

func TestSetWithTimeout(t *testing.T) {
	miner, err := GetMiner()

	assert.EqualValues(t, nil, err, "")
	if miner == nil {
		return
	}
	op, err := miner.SetWithTimeout("kk", "kiko", 10)
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}

func TestEvict(t *testing.T) {
	miner, err := GetMiner()

	assert.EqualValues(t, nil, err, "")
	if miner == nil {
		return
	}

	op, err := miner.Evict("kk")
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}

func TestClear(t *testing.T) {
	miner, err := GetMiner()

	assert.EqualValues(t, nil, err, "")
	if miner == nil {
		return
	}

	op, err := miner.Clear()
	assert.EqualValues(t, nil, err, "")
	assert.Equal(t, true, op, "")
}
