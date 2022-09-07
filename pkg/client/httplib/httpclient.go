package httplib

import (
	"github.com/valyala/fasthttp"
)

const (
	// 最大重定向次数
	defaultMaxRedirectsCount = 16
)

func New(opt *Options) (*HttpClient, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
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
	return c, nil
}

type HttpClient struct {
	client *fasthttp.Client
	opt    Options
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
