// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"

	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/go-redis/redis/v8"
)

type redisMetricKey string

const (
	keyRequestStart     redisMetricKey = "requestStart"
	keyPipeRequestStart redisMetricKey = "pipeRequestStart"
)

var _ redis.Hook = &metricHook{}

type metricHook struct {
	container *redisContainer
	logger    *log.NgoLogger
}

func newMetricHook(container *redisContainer) *metricHook {
	return &metricHook{
		container: container,
		logger:    log.WithField("name", container.opt.Name),
	}
}

func (h *metricHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *metricHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (h *metricHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *metricHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
