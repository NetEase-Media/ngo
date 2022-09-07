package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

func New(opt *Options) (*MemcacheProxy, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	c := memcache.New(opt.Addr...)
	c.Timeout = opt.Timeout
	c.MaxIdleConns = opt.MaxIdleConns
	p := &MemcacheProxy{
		base: c,
	}
	return p, nil
}

// MemcacheProxy memcache 三方包的包装器类
type MemcacheProxy struct {
	Opt  *Options
	base *memcache.Client
}
