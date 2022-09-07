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

func init() {

}

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
	// block 不调用 Completed，所以这里手动调用，进行统计
	// https://github.com/alibaba/sentinel-golang/blob/2af86608857884da543eabe7941b61d60edc51a4/core/base/slot_chain.go#L229
	c.OnCompleted(ctx)
}

func (c *MetricStatSlot) OnCompleted(ctx *base.EntryContext) {
}
