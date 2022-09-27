package pinpoint

import (
	"context"
	"net/http"

	"github.com/NetEase-Media/ngo/pkg/tracing"
	"github.com/NetEase-Media/ngo/pkg/util"
	"github.com/pinpoint-apm/pinpoint-go-agent"
)

var _ tracing.Tracer = new(Tracer)
var _ tracing.Span = new(Span)

var agentId string

func init() {
	tracing.Register(tracing.Pinpoint, New)
}

func SetAgentId(id string) {
	agentId = id
}

func New(opt *tracing.Options) (tracing.Tracer, error) {
	aid := opt.Pinpoint.AgentId
	if agentId != "" {
		aid = agentId
	}
	config, err := pinpoint.NewConfig(
		pinpoint.WithAppName(opt.Pinpoint.ApplicationName),
		pinpoint.WithAgentId(aid),
		pinpoint.WithCollectorHost(opt.Pinpoint.Collector.Host),
		pinpoint.WithCollectorAgentPort(opt.Pinpoint.Collector.AgentPort),
		pinpoint.WithCollectorSpanPort(opt.Pinpoint.Collector.SpanPort),
		pinpoint.WithCollectorStatPort(opt.Pinpoint.Collector.StatPort),
		pinpoint.WithSamplingRate(opt.Pinpoint.Sampling.Rate),
	)
	if err != nil {
		return nil, err
	}
	agent, err := pinpoint.NewAgent(config)
	if err != nil {
		return nil, err
	}
	return &Tracer{tracer: agent, enabled: opt.Enabled}, nil
}

type Tracer struct {
	tracer  *pinpoint.Agent
	enabled bool
}

func (t *Tracer) Enabled() bool {
	return t.enabled
}

func (t *Tracer) Type() string {
	return tracing.Pinpoint
}

func (t *Tracer) Inject(ctx context.Context, carrier interface{}) {
	if carrier == nil {
		return
	}
	span := tracing.SpanFromContext(ctx)
	if v, ok := carrier.(pinpoint.DistributedTracingContextWriter); ok {
		span.(*Span).span.Inject(v)
	}
}

func (t *Tracer) StartSpanFromCarrier(ctx context.Context, op string, carrier interface{}) (tracing.Span,
	context.Context) {
	if carrier != nil {
		switch carrier.(type) {
		case http.Header:
			s := t.tracer.NewSpanTracerWithReader(op, carrier.(http.Header))
			s = s.NewSpanEvent(op)
			span := &Span{span: s, root: true}
			return span, tracing.ContextWithSpan(ctx, span)
		}
	}
	return t.StartSpanFromContext(ctx, op)
}

func (t *Tracer) StartSpanFromContext(ctx context.Context, op string) (tracing.Span,
	context.Context) {
	span := tracing.SpanFromContext(ctx)
	var s pinpoint.Tracer
	if span != nil {
		s = span.(*Span).span.NewSpanEvent(op)
	} else {
		s = t.tracer.NewSpanTracer(op)
		s = s.NewSpanEvent(op)
		span = &Span{span: s, root: true}
		return span, tracing.ContextWithSpan(ctx, span)
	}
	span = &Span{span: s}
	return span, tracing.ContextWithSpan(ctx, span)
}

func (t *Tracer) SetSamplingRate(rate int) {
	t.tracer.SetSampling(pinpoint.Sampling{Rate: rate})
}

func (t *Tracer) Stop() {
	t.tracer.Shutdown()
}

type Span struct {
	root bool
	span pinpoint.Tracer
}

func (s *Span) SetTag(key string, value interface{}) tracing.Span {
	se := s.span.SpanEvent()
	switch key {
	case string(tracing.SpanType):
		switch util.Strval(value) {
		case "HTTP_SERVER":
			//se.SetServiceType(pinpoint.ServiceTypeGoApp)
		case "HTTP_CLIENT":
			se.SetServiceType(pinpoint.ServiceTypeGoHttpClient)
		case "REDIS":
			se.SetServiceType(8200)
		case "JDBC":
			se.SetServiceType(2101)
		}
	case string(tracing.PluginType):
	case string(tracing.RemoteAppName):
		se.SetDestination(util.Strval(value))
	case string(tracing.RemoteClusterName):
	case string(tracing.HttpServerRequestUrl):
	case string(tracing.HttpServerRequestMethod):
	case string(tracing.HttpServerRequestHost):
		s.span.Span().SetEndPoint(util.Strval(value))
	case string(tracing.HttpServerRequestPath):
		s.span.Span().SetRpcName(util.Strval(value))
	case string(tracing.HttpServerRequestSize):
	case string(tracing.HttpServerPeerHost):
		s.span.Span().SetRemoteAddress(util.Strval(value))
	case string(tracing.HttpServerPeerPort):
	case string(tracing.HttpServerResponseStatus):
		s.span.Span().Annotations().AppendInt(pinpoint.AnnotationHttpStatusCode, int32(value.(uint16)))
	case string(tracing.HttpServerResponseSize):
	case string(tracing.HttpClientRequestUrl):
		se.Annotations().AppendString(pinpoint.AnnotationHttpUrl, util.Strval(value))
	case string(tracing.HttpClientRequestMethod):
	case string(tracing.HttpClientRequestPath):
	case string(tracing.HttpClientRequestHost):
	case string(tracing.HttpClientResponseStatus):
		se.Annotations().AppendString(pinpoint.AnnotationHttpStatusCode, util.Strval(value))
	case string(tracing.DbUrl):
		se.SetEndPoint(util.Strval(value))
	case string(tracing.DbType):
	case string(tracing.DbSql):
		se.SetSQL(util.Strval(value))
	case string(tracing.RedisAddress):
		se.SetEndPoint(util.Strval(value))
		se.SetDestination(util.Strval(value))
	case string(tracing.RedisCmd):
	case string(tracing.RedisCode):
	case string(tracing.ExceptionName):
		se.SetError(value.(error))
	case string(tracing.ExceptionMessage):
	case string(tracing.ExceptionStacktrace):
	}
	return s
}

func (s *Span) Finish() {
	s.span.EndSpanEvent()
	if s.root {
		s.span.EndSpan()
	}
}

func (s *Span) GetTraceId() string {
	return s.span.TransactionId().String()
}
