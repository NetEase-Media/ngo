package tracing

import (
	"context"
)

var defaultTracer Tracer

func SetTracer(tracer Tracer) {
	defaultTracer = tracer
}

func Enabled() bool {
	return defaultTracer != nil && defaultTracer.Enabled()
}

func Type() string {
	return defaultTracer.Type()
}

func Inject(ctx context.Context, carrier interface{}) {
	defaultTracer.Inject(ctx, carrier)
}

func StartSpanFromContext(ctx context.Context, op string) (Span, context.Context) {
	return defaultTracer.StartSpanFromContext(ctx, op)
}

func StartSpanFromCarrier(ctx context.Context, op string, carrier interface{}) (Span, context.Context) {
	return defaultTracer.StartSpanFromCarrier(ctx, op, carrier)
}

func SetSamplingRate(rate int) {
	defaultTracer.SetSamplingRate(rate)
}

func Stop() {
	defaultTracer.Stop()
}
