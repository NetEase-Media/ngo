package opentracing

import "github.com/valyala/fasthttp"

type FasthttpCarrier struct {
	header *fasthttp.RequestHeader
}

func NewFasthttpCarrier(header *fasthttp.RequestHeader) *FasthttpCarrier {
	return &FasthttpCarrier{header: header}
}

// TextMapWriter
func (carrier *FasthttpCarrier) Set(key, val string) {
	carrier.header.Set(key, val)
}

// TextMapReader
func (carrier *FasthttpCarrier) ForeachKey(handler func(key, val string) error) (err error) {
	carrier.header.VisitAll(func(k, v []byte) {
		if err = handler(string(k), string(v)); err != nil {
			return
		}
	})
	return err
}
