package redis

import (
	"context"
	"runtime"

	"github.com/NetEase-Media/ngo/pkg/tracing"

	"github.com/go-redis/redis/extra/rediscmd"
	"github.com/go-redis/redis/v8"
)

type tracingHook struct {
	container *RedisContainer
}

func newTracingHook(container *RedisContainer) *tracingHook {
	return &tracingHook{container: container}
}

func (th *tracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !tracing.Enabled() {
		return ctx, nil
	}
	span, c := tracing.StartSpanFromContext(ctx, "redis:"+cmd.FullName())
	tracing.SpanKind.Set(span, "client")
	tracing.SpanType.Set(span, "REDIS")
	tracing.PluginType.Set(span, "go-redis")
	tracing.RedisAddress.Set(span, th.container.Opt.Addr)
	tracing.RemoteAppName.Set(span, "redis")
	return c, nil
}

func (th *tracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if !tracing.Enabled() {
		return nil
	}
	span := tracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	if err := cmd.Err(); err != nil {
		recordError(ctx, span, err)
		tracing.RedisCode.Set(span, "-1")
	} else {
		tracing.RedisCode.Set(span, "200")
	}
	span.Finish()
	return nil
}

func (th *tracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !tracing.Enabled() {
		return ctx, nil
	}
	summary, cmdsString := rediscmd.CmdsString(cmds)
	span, c := tracing.StartSpanFromContext(ctx, "redis-pipeline:"+summary)
	tracing.SpanKind.Set(span, "client")
	tracing.SpanType.Set(span, "REDIS")
	tracing.PluginType.Set(span, "go-redis")
	tracing.RedisAddress.Set(span, th.container.Opt.Addr)
	tracing.RedisCmd.Set(span, cmdsString)
	tracing.RemoteAppName.Set(span, "redis")

	return c, nil
}

func (th *tracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if !tracing.Enabled() {
		return nil
	}
	span := tracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	resultCode := "200"
	if err := cmds[0].Err(); err != nil {
		recordError(ctx, span, err)
		resultCode = "-1"
	}
	tracing.RedisCode.Set(span, resultCode)
	span.Finish()
	return nil
}

func recordError(ctx context.Context, span tracing.Span, err error) {
	if err != redis.Nil {
		stack := make([]byte, 2048)
		runtime.Stack(stack, false)
		tracing.ExceptionName.Set(span, err)
		tracing.ExceptionMessage.Set(span, err.Error())
		tracing.ExceptionStacktrace.Set(span, string(stack))
	}
}
