package redis

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	opts := &Options{
		Name:     "client1",
		Addr:     []string{"127.0.0.1:2379"},
		ConnType: "client",
	}

	client, err := New(opts)
	assert.Nil(t, err)
	assert.Equal(t, []string{"127.0.0.1:2379"}, client.Opt.Addr)
	client.Close()
}

func generateKey() string {
	t := time.Now().Unix()
	return "test_key:" + strconv.FormatInt(t, 10) + "_" + strconv.FormatInt(int64(rand.Int()), 10)
}
