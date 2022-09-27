package tracing

import (
	"context"
)

type ContextKey string

const (
	contextKey ContextKey = "ngo.spanTracer"
)

const (
	SpanKind   = StringTagName("span.kind")
	SpanType   = StringTagName("span.type")
	PluginType = StringTagName("span.plugin.type")

	RemoteAppName     = StringTagName("remote.app.name")
	RemoteClusterName = StringTagName("remote.cluster.name")

	HttpServerRequestUrl     = StringTagName("http.server.request.url")
	HttpServerRequestMethod  = StringTagName("http.server.request.method")
	HttpServerRequestHost    = StringTagName("http.server.request.host")
	HttpServerRequestPath    = StringTagName("http.server.request.path")
	HttpServerRequestSize    = StringTagName("http.server.request.size")
	HttpServerPeerHost       = StringTagName("http.server.peer.host")
	HttpServerPeerPort       = StringTagName("http.server.peer.port")
	HttpServerResponseStatus = StringTagName("http.server.response.status")
	HttpServerResponseSize   = StringTagName("http.server.response.size")

	HttpClientRequestUrl     = StringTagName("http.client.request.url")
	HttpClientRequestMethod  = StringTagName("http.client.request.method")
	HttpClientRequestHost    = StringTagName("http.client.request.host")
	HttpClientRequestPath    = StringTagName("http.client.request.path")
	HttpClientResponseStatus = StringTagName("http.client.response.status")

	DbUrl  = StringTagName("db.url")
	DbType = StringTagName("db.type")
	DbSql  = StringTagName("db.sql")

	RedisAddress = StringTagName("redis.address")
	RedisCmd     = StringTagName("redis.cmd")
	RedisCode    = StringTagName("redis.code")

	ExceptionName       = StringTagName("exception.name")
	ExceptionMessage    = StringTagName("exception.message")
	ExceptionStacktrace = StringTagName("exception.stacktrace")
)

type Tracer interface {
	Enabled() bool
	Type() string
	StartSpanFromCarrier(ctx context.Context, op string, carrier interface{}) (Span, context.Context)
	StartSpanFromContext(ctx context.Context, op string) (Span, context.Context)
	Inject(ctx context.Context, carrier interface{})
	SetSamplingRate(rate int)
	Stop()
}

type Span interface {
	SetTag(key string, value interface{}) Span
	Finish()
	GetTraceId() string
}

func SpanFromContext(ctx context.Context) Span {
	v := ctx.Value(contextKey)
	if v != nil {
		return v.(Span)
	}
	return nil
}

func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, contextKey, span)
}

func GetTraceId(ctx context.Context) string {
	span := SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	return span.GetTraceId()
}

// StringTagName is a common tag name to be set to a string value
type StringTagName string

// Set adds a string tag to the `span`
func (tag StringTagName) Set(span Span, value interface{}) {
	span.SetTag(string(tag), value)
}
