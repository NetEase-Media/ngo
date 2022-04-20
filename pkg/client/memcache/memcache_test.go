package memcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	o := Options{
		Name: "m1",
		Addr: []string{"127.0.0.1:11312"},
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "encounter error.")
	assert.NotEqual(t, nil, c, "Init Client Failed")

	data := make(map[string]string)
	data["Halo"] = "World"
	data["Halo1"] = "World1"

	// 插入
	for k, v := range data {
		err := c.Set(k, v)
		assert.Equal(t, nil, err, "err.")
	}

	// 查询
	for k := range data {
		d, err := c.Get(k)
		expect := data[k]
		assert.Equal(t, nil, err, "err.")
		assert.Equal(t, expect, d, "wrong data.")
	}

	// 批量查询
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	m, err := c.MGet(keys)
	assert.Equal(t, nil, err, "err.")
	assert.Equal(t, data, m, "mget wrong data.")

	// 删除
	for k := range data {
		err := c.Delete(k)
		assert.Equal(t, nil, err, "err.")
	}

	// 校验删除结果
	m, err = c.MGet(keys)
	assert.Equal(t, nil, err, "err.")
	assert.Equal(t, 0, len(m), "delete failed")
}
