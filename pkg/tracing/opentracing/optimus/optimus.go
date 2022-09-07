package optimus

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/NetEase-Media/ngo/jaeger-client-go"
	jaegercfg "github.com/NetEase-Media/ngo/jaeger-client-go/config"
	"github.com/NetEase-Media/ngo/pkg/env"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/tracing"
	ot "github.com/NetEase-Media/ngo/pkg/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/valyala/fasthttp"
)

var (
	_ tracing.Tracer = new(Tracer)
	_ tracing.Span   = new(Span)
)

func init() {
	tracing.Register(tracing.Optimus, New)
}

func New(opt *tracing.Options) (tracing.Tracer, error) {
	config := opt.Optimus

	var address string
	if config.UdpHost == "" {
		return nil, errors.New("empty UdpHost")
	}
	address = config.UdpHost
	address += ":"
	if len(config.UdpPort) > 0 {
		address += config.UdpPort
	} else {
		address += "6831"
	}
	udpTransport, err := jaeger.NewUDPTransport(address, 0)
	if err != nil {
		log.Errorf("create udp transport error:%v", err)
		return nil, err
	}
	cfg := jaegercfg.Configuration{
		ServiceName: config.ServiceName,
	}
	logWrapper := newNgoLoggerWrapper(log.DefaultLogger().(*log.NgoLogger).WithOptions(zap.AddCallerSkip(1)))
	ngoPropagator := NewNgoB3HTTPHeaderPropagator()

	tracer, gCloser, err := cfg.NewTracer(
		jaegercfg.Logger(logWrapper),
		jaegercfg.Metrics(metrics.NullFactory),
		jaegercfg.Injector(opentracing.HTTPHeaders, ngoPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, ngoPropagator),
		jaegercfg.Sampler(jaeger.NewConstSampler(true)),
		jaegercfg.Reporter(jaeger.NewRemoteReporter(udpTransport)),
		jaegercfg.Gen128Bit(true),
		jaegercfg.ZipkinSharedRPCSpan(true),
		jaegercfg.Tag("local.app.name", env.GetAppName()),
		jaegercfg.Tag("local.app.clusterName", env.GetClusterName()),
	)
	if err != nil {
		logWrapper.Errorf("Could not initialize optimus tracer: %s", err.Error())
		return nil, err
	}
	return &Tracer{tracer: tracer, gCloser: gCloser, enabled: opt.Enabled}, nil
}

type Tracer struct {
	tracer  opentracing.Tracer
	gCloser io.Closer
	enabled bool
}

func (t *Tracer) Enabled() bool {
	return t.enabled
}

func (t *Tracer) Type() string {
	return tracing.Optimus
}

func (t *Tracer) Inject(ctx context.Context, carrier interface{}) {
	if carrier == nil {
		return
	}
	span := tracing.SpanFromContext(ctx)
	switch carrier.(type) {
	case *fasthttp.RequestHeader:
		carrier := ot.NewFasthttpCarrier(carrier.(*fasthttp.RequestHeader))
		t.tracer.Inject(span.(*Span).span.Context(), opentracing.HTTPHeaders, carrier)
	}

}

func (t *Tracer) StartSpanFromCarrier(ctx context.Context, op string, carrier interface{}) (tracing.Span, context.Context) {
	if carrier != nil {
		switch carrier.(type) {
		case http.Header:
			spanContext, err := t.tracer.Extract(opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(carrier.(http.Header)))
			var s opentracing.Span
			if err != nil {
				log.Errorf("extract error, %s", err)
				s = t.tracer.StartSpan(op, ext.SpanKindRPCServer)
				s.SetTag("root", true)
			} else {
				s = t.tracer.StartSpan(op, opentracing.ChildOf(spanContext), ext.SpanKindRPCServer)
			}
			span := &Span{span: s}
			return span, tracing.ContextWithSpan(ctx, span)
		}
	}

	span, c := t.StartSpanFromContext(ctx, op)
	span.SetTag("root", true)
	return span, c
}

func (t *Tracer) StartSpanFromContext(ctx context.Context, op string) (tracing.Span,
	context.Context) {
	parentSpan := tracing.SpanFromContext(ctx)
	var s opentracing.Span
	if parentSpan != nil {
		s = t.tracer.StartSpan(op, opentracing.ChildOf(parentSpan.(*Span).span.Context()))
	} else {
		s = t.tracer.StartSpan(op)
	}
	span := &Span{span: s}
	return span, tracing.ContextWithSpan(ctx, span)
}

func (t *Tracer) SetSamplingRate(rate int) {

}

func (t *Tracer) Stop() {
	t.gCloser.Close()
}

type Span struct {
	span opentracing.Span
}

func (s *Span) SetTag(key string, value interface{}) tracing.Span {
	s.span.SetTag(key, value)
	return s
}

func (s *Span) Finish() {
	s.span.Finish()
}

func (s *Span) GetTraceId() string {
	return s.span.Context().(jaeger.SpanContext).TraceID().String()
}
