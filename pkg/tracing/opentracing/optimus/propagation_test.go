package optimus

import (
	"testing"

	"github.com/NetEase-Media/ngo/jaeger-client-go"
	"github.com/stretchr/testify/assert"
)

type testWriter struct {
	container map[string]string
}

func newTestWriter() *testWriter {
	return &testWriter{container: make(map[string]string)}
}

func (tw *testWriter) Set(key, val string) {
	tw.container[key] = val
}

func (tw *testWriter) ForeachKey(handler func(key, val string) error) error {
	for k := range tw.container {
		handler(k, tw.container[k])
	}
	return nil
}

func TestInjectAndExtract(t *testing.T) {
	p := NewNgoB3HTTPHeaderPropagator(BaggagePrefix("baggage-"))
	assert.Equal(t, "baggage-", p.baggagePrefix)

	spanContext := jaeger.NewSpanContext("hello_001", "0.1", "0", false, nil)
	tw := newTestWriter()
	err := p.Inject(spanContext, tw)
	if err != nil {
		t.Fatalf("error:%v", err)
		return
	}
	assert.Equal(t, "hello_001", tw.container["x-b3-traceid"])
	assert.Equal(t, "0.1", tw.container["x-b3-spanid"])
	assert.Equal(t, "0", tw.container["x-b3-parentspanid"])
	assert.Equal(t, "0", tw.container["x-b3-sampled"])
	assert.Equal(t, "hello_001", tw.container["x-b3-traceid"])

	sc, err := p.Extract(tw)
	if err != nil {
		t.Fatalf("error:%v", err)
	}

	assert.Equal(t, "hello_001", sc.TraceID().String())
	assert.Equal(t, "0", sc.ParentID().String())
	assert.Equal(t, "0.1", sc.SpanID().String())
	assert.Equal(t, false, sc.IsSampled())

}
