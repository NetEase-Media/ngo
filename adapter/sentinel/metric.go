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

package sentinel

import (
	"github.com/alibaba/sentinel-golang/core/base"
)

const (
	StatSlotOrder = 9999

	metricsKey = "sentinelRequestKey"
)

var (
	NssMetricSlot = &MetricStatSlot{}
)

func init() {}

// MetricStatSlot records metrics for circuit breaker on invocation completed.
// MetricStatSlot must be filled into slot chain if circuit breaker is alive.
type MetricStatSlot struct {
	base.StatPrepareSlot
}

func (s *MetricStatSlot) Order() uint32 {
	return StatSlotOrder
}

func (c *MetricStatSlot) Prepare(ctx *base.EntryContext) {

}

func (c *MetricStatSlot) OnEntryPassed(ctx *base.EntryContext) {
}

func (c *MetricStatSlot) OnEntryBlocked(ctx *base.EntryContext, err *base.BlockError) {

}

func (c *MetricStatSlot) OnCompleted(ctx *base.EntryContext) {
}
