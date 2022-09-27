package sentinel

import (
	"testing"

	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/stretchr/testify/assert"
)

func TestSentinel(t *testing.T) {
	opt := Options{
		FlowRules: []*flow.Rule{
			{
				Resource:               "some-test",
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              1,
				StatIntervalInMs:       10000,
			},
		},
	}
	Init(&opt)
	var succ, fail int
	for i := 0; i < 2; i++ {
		e, b := Entry("some-test")
		if b != nil {
			fail++
		} else {
			succ++
			e.Exit()
		}
	}
	assert.Equal(t, succ, 1)
	assert.Equal(t, fail, 1)
}
