package httplib

import (
	"crypto/tls"
	"time"
)

type Options struct {
	// User-Agent header，如果为空会设置默认header
	Name string

	// 如果为true，即使Name为空也不设置默认User-Agent
	NoDefaultUserAgentHeader bool

	// TLS配置，需要拆分
	TLSConfig *tls.Config

	// 每个host的最大连接数
	MaxConnsPerHost int

	// 空闲的keep-alive连接最大关闭时间
	MaxIdleConnDuration time.Duration

	// keep-alive连接最大关闭时间
	MaxConnDuration time.Duration

	// 最大重试次数
	MaxIdemponentCallAttempts int

	// 每个连接的读缓存大小，会限制最大header长度
	ReadBufferSize int

	// 每个连接的写缓存大小
	WriteBufferSize int

	// 最大的回复读取时间
	ReadTimeout time.Duration

	// 最大的写请求事件
	WriteTimeout time.Duration

	// 最大的回复body大小
	MaxResponseBodySize int

	// 等待空闲连接的最大时间。默认情况不等待，如果没有空闲连接返回ErrNoFreeConns错误。
	MaxConnWaitTimeout time.Duration
}

func NewDefaultOptions() *Options {
	return &Options{
		Name:                     "",
		NoDefaultUserAgentHeader: false,
		TLSConfig:                nil,
		MaxConnsPerHost:          512,
		MaxIdleConnDuration:      time.Second * 10,
		//MaxConnDuration:
		MaxIdemponentCallAttempts: 5,
		ReadBufferSize:            4096,
		WriteBufferSize:           4096,
		ReadTimeout:               time.Second * 60,
		WriteTimeout:              time.Second * 60,
		//MaxResponseBodySize:
		//MaxConnWaitTimeout:
	}
}

func checkOptions(opt *Options) error {
	return nil
}
