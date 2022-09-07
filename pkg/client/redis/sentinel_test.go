package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSentinelClient(t *testing.T) {
	// 有外部依赖，先忽略
	t.Skip()
	ctx := context.Background()
	client := NewSentinelClient(&Options{
		Name:         "sharedsentinel01",
		ConnType:     "sentinel",
		MasterNames:  []string{"recncr6510"},
		Addr:         []string{"127.0.0.1:26381"},
		PoolSize:     50,
		Password:     "pushsentinel",
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})
	defer client.Close()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "test:" + strconv.Itoa(i)
			client.Set(ctx, key, "value"+strconv.Itoa(i), time.Second*10)
			stringCmd := client.Get(ctx, key)
			fmt.Println(stringCmd.Val())
		}(i)
	}
	wg.Wait()
}

func BenchmarkSentinelClient_Set(b *testing.B) {
	// 有外部依赖，先忽略
	b.Skip()
	ctx := context.Background()
	client := NewSentinelClient(&Options{
		Name:         "sharedsentinel01",
		ConnType:     "sentinel",
		MasterNames:  []string{"recncr6510"},
		Addr:         []string{"127.0.0.1:26381"},
		PoolSize:     50,
		Password:     "pushsentinel",
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := client.Set(ctx, "test:"+strconv.Itoa(i), "value"+strconv.Itoa(i), time.Second*10)
		assert.NoError(b, cmd.Err())
	}
	b.StopTimer()
}
