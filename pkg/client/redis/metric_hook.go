package redis

//
//import (
//	"context"
//
//	nssredis "g.hz.netease.com/agent/nss-go-agent/collector/collectors/redis"
//	"github.com/NetEase-Media/ngo/pkg/log"
//	"github.com/NetEase-Media/ngo/pkg/metrics"
//	collectors "github.com/NetEase-Media/ngo/pkg/metrics/colloctors"
//	"github.com/go-redis/redis/v8"
//)
//
//type redisMetricKey string
//
//const (
//	keyRequestStart     redisMetricKey = "requestStart"
//	keyPipeRequestStart redisMetricKey = "pipeRequestStart"
//)
//
//var _ redis.Hook = &metricHook{}
//
//type metricHook struct {
//	container *RedisContainer
//	logger    log.Logger
//}
//
//func newMetricHook(container *RedisContainer) *metricHook {
//	return &metricHook{
//		container: container,
//		logger:    log.WithField("name", container.Opt.Name),
//	}
//}
//
//func (h *metricHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
//	if !metrics.IsMetricsEnabled() {
//		return ctx, nil
//	}
//	host := h.container.Opt.Addr[0]
//	stats := collectors.RedisCollector().OnStart(host, cmd.Name())
//	ctx = context.WithValue(ctx, keyRequestStart, stats)
//	return ctx, nil
//}
//
//func (h *metricHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
//	if !metrics.IsMetricsEnabled() {
//		return nil
//	}
//
//	// 万一有人改写了数据，这里只能打日志再退出
//	v := ctx.Value(keyRequestStart)
//	if v == nil {
//		h.logger.Errorf("can not get value from %s", keyRequestStart)
//		return nil
//	}
//
//	stats, ok := v.(*nssredis.StatsHolder)
//
//	if !ok {
//		h.logger.Errorf("convert to redis.Stats failed from %s", keyRequestStart)
//		return nil
//	}
//
//	// 忽略key不存在情况
//	if cmd.Err() != nil && cmd.Err() != redis.Nil {
//		collectors.RedisCollector().OnError(stats, cmd.Err())
//	}
//
//	collectors.RedisCollector().OnComplete(stats, cmd.Err() == nil)
//
//	return nil
//}
//
//func (h *metricHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
//	if !metrics.IsMetricsEnabled() {
//		return ctx, nil
//	}
//	statses := make(map[redis.Cmder]*nssredis.StatsHolder, len(cmds))
//	for i := range cmds {
//		host := h.container.Opt.Addr[0]
//		statses[cmds[i]] = collectors.RedisCollector().OnStart(host, cmds[i].Name())
//	}
//	ctx = context.WithValue(ctx, keyPipeRequestStart, statses)
//	return ctx, nil
//}
//
//func (h *metricHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
//	if !metrics.IsMetricsEnabled() {
//		return nil
//	}
//
//	// 万一有人改写了数据，这里只能打日志再退出
//	v := ctx.Value(keyPipeRequestStart)
//	if v == nil {
//		h.logger.Errorf("can not get value from %s", keyPipeRequestStart)
//		return nil
//	}
//
//	statses, ok := v.(map[redis.Cmder]*nssredis.StatsHolder)
//
//	if !ok {
//		h.logger.Errorf("convert to redis.Stats failed from %s", keyPipeRequestStart)
//		return nil
//	}
//
//	for i := range cmds {
//		// 忽略key不存在情况
//		if cmds[i].Err() != nil && cmds[i].Err() != redis.Nil {
//			collectors.RedisCollector().OnError(statses[cmds[i]], cmds[i].Err())
//		}
//
//		collectors.RedisCollector().OnComplete(statses[cmds[i]], cmds[i].Err() == nil)
//	}
//	return nil
//}
