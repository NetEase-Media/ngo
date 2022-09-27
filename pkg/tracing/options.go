package tracing

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	Jaeger   = "jaeger"
	Zipkin   = "zipkin"
	Pinpoint = "pinpoint"
	Optimus  = "optimus"
)

type Options struct {
	Enabled  bool
	Type     string
	Pinpoint PinpointConfig
	Optimus  OptimusConfig
	Baggage  map[string]string
}

type PinpointConfig struct {
	ApplicationName string
	ApplicationType int32
	AgentId         string
	ConfigFilePath  string

	Collector struct {
		Host      string
		AgentPort int
		SpanPort  int
		StatPort  int
	}

	LogLevel logrus.Level

	Sampling struct {
		Rate               int
		NewThroughput      int
		ContinueThroughput int
	}

	Stat struct {
		CollectInterval int
		BatchCount      int
	}
}

type OptimusConfig struct {
	ServiceName string
	UdpHost     string
	UdpPort     string
}

func NewDefaultOptions() *Options {
	return &Options{
		Enabled: false,
	}
}

func checkOptions(opt *Options) error {
	if opt.Type == Pinpoint {
		hostname, _ := os.Hostname()
		opt.Pinpoint.AgentId = hostname
	}
	return nil
}
