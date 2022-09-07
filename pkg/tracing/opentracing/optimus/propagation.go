package optimus

import (
	"strings"

	"github.com/NetEase-Media/ngo/jaeger-client-go"
	"github.com/opentracing/opentracing-go"
)

// Option is a function that sets an option on Propagator
type Option func(propagator *Propagator)

// BaggagePrefix is a function that sets baggage prefix on Propagator
func BaggagePrefix(prefix string) Option {
	return func(propagator *Propagator) {
		propagator.baggagePrefix = prefix
	}
}

// Propagator is an Injector and Extractor
type Propagator struct {
	baggagePrefix string
}

func NewNgoB3HTTPHeaderPropagator(opts ...Option) Propagator {
	p := Propagator{baggagePrefix: "baggage-"}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// Inject conforms to the Injector interface for decoding Zipkin HTTP B3 headers
func (p Propagator) Inject(
	sc jaeger.SpanContext,
	abstractCarrier interface{},
) error {
	textMapWriter, ok := abstractCarrier.(opentracing.TextMapWriter)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}

	textMapWriter.Set("x-b3-traceid", sc.TraceID().String())
	if len(sc.ParentID()) > 0 {
		textMapWriter.Set("x-b3-parentspanid", sc.ParentID().String())
	}
	textMapWriter.Set("x-b3-spanid", sc.SpanID().String())
	if sc.IsSampled() {
		textMapWriter.Set("x-b3-sampled", "1")
	} else {
		textMapWriter.Set("x-b3-sampled", "0")
	}
	sc.ForeachBaggageItem(func(k, v string) bool {
		textMapWriter.Set(p.baggagePrefix+k, v)
		return true
	})
	return nil
}

// Extract conforms to the Extractor interface for encoding Zipkin HTTP B3 headers
func (p Propagator) Extract(abstractCarrier interface{}) (jaeger.SpanContext, error) {
	textMapReader, ok := abstractCarrier.(opentracing.TextMapReader)
	if !ok {
		return jaeger.SpanContext{}, opentracing.ErrInvalidCarrier
	}
	var traceID jaeger.TraceID
	var spanID jaeger.SpanID
	var parentID jaeger.SpanID
	sampled := false
	baggage := make(map[string]string)
	err := textMapReader.ForeachKey(func(rawKey, value string) error {
		key := strings.ToLower(rawKey)
		var err error
		if len(value) <= 0 {
			return err
		}
		if key == "x-b3-traceid" {
			traceID, err = jaeger.TraceIDFromString(value)
		} else if key == "x-b3-parentspanid" {
			parentID, err = jaeger.SpanIDFromString(value)
		} else if key == "x-b3-spanid" {
			spanID, err = jaeger.SpanIDFromString(value)
		} else if key == "x-b3-sampled" && (value == "1" || value == "true") {
			sampled = true
		} else if strings.HasPrefix(key, p.baggagePrefix) {
			baggage[key[len(p.baggagePrefix):]] = value
		}
		return err
	})

	if err != nil {
		return jaeger.SpanContext{}, err
	}
	if !traceID.IsValid() {
		return jaeger.SpanContext{}, opentracing.ErrSpanContextNotFound
	}
	return jaeger.NewSpanContext(traceID, spanID, parentID, sampled, baggage), nil
}
