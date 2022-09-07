package sentinel

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/pkg/log"

	"github.com/alibaba/sentinel-golang/core/flow"

	"github.com/stretchr/testify/assert"

	"github.com/NetEase-Media/ngo/pkg/metrics"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
)

func TestMetrics(t *testing.T) {
	metrics.Init("ngo-demo", "ngo-demo-docker-cm_test")

	err := Init(&Options{
		CircuitBreakerRules: []*circuitbreaker.Rule{
			// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
			{
				Resource:         "abc",
				Strategy:         circuitbreaker.ErrorCount,
				RetryTimeoutMs:   3000,
				MinRequestAmount: 1,
				StatIntervalMs:   5000,
				Threshold:        1,
			},
		},
		FlowRules: []*flow.Rule{
			{
				Resource:               "abc",
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              100,
				StatIntervalInMs:       1000,
			},
			{
				Resource:               "abc",
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              1,
				StatIntervalInMs:       1000,
			},
		},
	})
	assert.Nil(t, err)

	for i := 0; i < 3; i++ {
		e, b := Entry("abc")
		if b != nil {
			log.Errorf("%s", b)
			// g1 blocked
			time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
		} else {
			if rand.Uint64()%20 > 9 {
				// Record current invocation as error.
				TraceError(e, errors.New("biz error"))
			}
			// g1 passed
			time.Sleep(time.Duration(rand.Uint64()%80+10) * time.Millisecond)
			e.Exit()
		}
	}
}
