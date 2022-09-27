package ratelimiter

import (
	"os"
	"testing"

	"github.com/NetEase-Media/ngo/pkg/sentinel"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func TestMain(m *testing.M) {
	setupTest()
	ret := m.Run()
	tearDownTest()
	os.Exit(ret)
}

func setupTest() {
	sentinel.Init(&sentinel.Options{
		FlowRules: []*flow.Rule{
			{
				Resource:               "abc",
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              1,
				StatIntervalInMs:       1000,
			},
			{
				Resource:               "def",
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              1,
				StatIntervalInMs:       1000,
			},
		},
	})
}

func tearDownTest() {

}
