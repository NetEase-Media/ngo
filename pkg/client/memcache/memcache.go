package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

// Get 根据Key获取缓存数值
func (m *MemcacheProxy) Get(key string) (string, error) {
	item, err := m.base.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

// MGet 获取多个数值
func (m *MemcacheProxy) MGet(keys []string) (map[string]string, error) {
	rets, err := m.base.GetMulti(keys)
	if err != nil {
		return nil, err
	}

	if len(rets) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(rets))
	for _, v := range rets {
		r[v.Key] = string(v.Value)
	}
	return r, nil
}

// Set 设置缓存
func (m *MemcacheProxy) Set(key string, value string) error {
	item := memcache.Item{
		Key:   key,
		Value: []byte(value),
	}
	err := m.base.Set(&item)
	return err
}

// SetWithExpire 设置缓存，并且添加超时
// expire 以s为单位
func (m *MemcacheProxy) SetWithExpire(key string, value string, expire int) error {
	item := memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(expire),
	}
	err := m.base.Set(&item)
	return err
}

// Delete 删除操作
func (m *MemcacheProxy) Delete(key string) error {
	err := m.base.Delete(key)
	return err
}
