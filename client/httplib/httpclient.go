// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httplib

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	defaultHttpClient *HttpClient
)

const (
	// 最大重定向次数
	defaultMaxRedirectsCount = 16
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
	return &Options{}
}

type HttpClient struct {
	client *fasthttp.Client
	opt    Options
}

func New(opt *Options) *HttpClient {
	client := &fasthttp.Client{
		Name:                      opt.Name,
		NoDefaultUserAgentHeader:  opt.NoDefaultUserAgentHeader,
		TLSConfig:                 opt.TLSConfig, // TODO: TLS系列配置需要另外分离
		MaxConnsPerHost:           opt.MaxConnsPerHost,
		MaxIdleConnDuration:       opt.MaxIdleConnDuration,
		MaxConnDuration:           opt.MaxConnDuration,
		MaxIdemponentCallAttempts: opt.MaxIdemponentCallAttempts,
		ReadBufferSize:            opt.ReadBufferSize,
		WriteBufferSize:           opt.WriteBufferSize,
		ReadTimeout:               opt.ReadTimeout,
		WriteTimeout:              opt.WriteTimeout,
		MaxResponseBodySize:       opt.MaxResponseBodySize,
		MaxConnWaitTimeout:        opt.MaxConnWaitTimeout,
	}
	c := &HttpClient{
		client: client,
		opt:    *opt,
	}

	return c
}

// InitialOptions 返回初始化时所使用的配置选项
func InitialOptions() Options {
	return defaultHttpClient.opt
}

func Init(opt *Options) {
	if defaultHttpClient != nil {
		panic("duplicated init http client")
	}

	defaultHttpClient = New(opt)
}

// Get 调用默认http客户端的GET方法
func Get(url string) *DataFlow {
	return defaultHttpClient.Get(url)
}

// Post 调用默认http客户端的POST方法
func Post(url string) *DataFlow {
	return defaultHttpClient.Post(url)
}

// Put 调用默认http客户端的PUT方法
func Put(url string) *DataFlow {
	return defaultHttpClient.Put(url)
}

// Delete 调用默认http客户端的DELETE方法
func Delete(url string) *DataFlow {
	return defaultHttpClient.Delete(url)
}

// Patch 调用默认http客户端的PATCH方法
func Patch(url string) *DataFlow {
	return defaultHttpClient.Patch(url)
}

func Close() {
	defaultHttpClient.Close()
}

// Get 调用http客户端的GET方法
func (c *HttpClient) Get(url string) *DataFlow {
	df := NewDataFlow(c)
	return df.newMethod(fasthttp.MethodGet, url)
}

// Post 调用http客户端的POST方法
func (c *HttpClient) Post(url string) *DataFlow {
	df := NewDataFlow(c)
	return df.newMethod(fasthttp.MethodPost, url)
}

// Put 调用http客户端的PUT方法
func (c *HttpClient) Put(url string) *DataFlow {
	df := NewDataFlow(c)
	return df.newMethod(fasthttp.MethodPut, url)
}

// Delete 调用http客户端的DELETE方法
func (c *HttpClient) Delete(url string) *DataFlow {
	df := NewDataFlow(c)
	return df.newMethod(fasthttp.MethodDelete, url)
}

// Patch 调用http客户端的PATCH方法
func (c *HttpClient) Patch(url string) *DataFlow {
	df := NewDataFlow(c)
	return df.newMethod(fasthttp.MethodPatch, url)
}

func (c *HttpClient) Close() {
	c.client.CloseIdleConnections()
}
