package redis

import (
	"context"
	"errors"
	"unsafe"

	"github.com/go-redis/redis/v8"
)

func NewShardedPipeline(ctx context.Context, fn func(context.Context, []redis.Cmder) error) redis.Pipeliner {
	pipe := &redis.Pipeline{}
	c := (*context.Context)(unsafe.Pointer(uintptr(unsafe.Pointer(pipe)) +
		unsafe.Sizeof(func(ctx context.Context, cmd redis.Cmder) error { return nil }) +
		unsafe.Sizeof(ctx)))
	*c = ctx
	e := (*func(context.Context, []redis.Cmder) error)(unsafe.Pointer(uintptr(unsafe.Pointer(pipe)) +
		unsafe.Sizeof(func(ctx context.Context, cmd redis.Cmder) error { return nil }) +
		unsafe.Sizeof(ctx) +
		unsafe.Sizeof(func(context.Context, []redis.Cmder) error { return nil })))
	*e = fn
	pipeInit(pipe)
	return pipe
}

//go:linkname pipeInit github.com/go-redis/redis/v8.(*Pipeline).init
func pipeInit(*redis.Pipeline) int

//go:linkname clientProcessPipeline github.com/go-redis/redis/v8.(*Client).processPipeline
func clientProcessPipeline(*redis.Client, context.Context, []redis.Cmder) error

//go:linkname clusterClientProcessPipeline github.com/go-redis/redis/v8.(*ClusterClient).processPipeline
func clusterClientProcessPipeline(*redis.ClusterClient, context.Context, []redis.Cmder) error

func processPipeline(ctx context.Context, client interface{}, cmds []redis.Cmder) error {
	rc := client.(*RedisContainer)
	switch rc.Redis.(type) {
	case *redis.Client:
		return clientProcessPipeline(rc.Redis.(*redis.Client), ctx, cmds)
	case *redis.ClusterClient:
		return clusterClientProcessPipeline(rc.Redis.(*redis.ClusterClient), ctx, cmds)
	default:
		return errors.New("client must be type of redis.Client or redis.ClusterClient")
	}
}
