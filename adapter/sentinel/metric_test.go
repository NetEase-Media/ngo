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
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"

	"github.com/alibaba/sentinel-golang/core/flow"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
)

func TestMetrics(t *testing.T) {

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
