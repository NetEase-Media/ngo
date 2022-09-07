package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	NAME = "ngo-test"
)

func TestInit_InitProcess(t *testing.T) {
	opts := NewDefaultOptions()
	opts.Name = NAME
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	k, err := New(opts)
	assert.Equal(t, nil, err)
	consumer := k.Consumer
	assert.Equal(t, "ngo", consumer.opt.Consumer.Group)
	producer := k.Producer
	assert.Equal(t, time.Second*10, producer.opt.Producer.Timeout)
}
