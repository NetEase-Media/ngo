package httplib

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/pkg/util"
)

func BenchmarkHTTP_1Conn_1Delay(b *testing.B) {
	testRunBench(b, 1, 1)
}

func BenchmarkHTTP_1Conn_5Delay(b *testing.B) {
	testRunBench(b, 1, 5)
}

func BenchmarkHTTP_1Conn_50Delay(b *testing.B) {
	testRunBench(b, 1, 50)
}

func BenchmarkHTTP_5Conn_1Delay(b *testing.B) {
	testRunBench(b, 5, 1)
}

func BenchmarkHTTP_5Conn_5Delay(b *testing.B) {
	testRunBench(b, 5, 5)
}

func BenchmarkHTTP_5Conn_50Delay(b *testing.B) {
	testRunBench(b, 5, 50)
}

func BenchmarkHTTP_100Conn_1Delay(b *testing.B) {
	testRunBench(b, 100, 1)
}

func BenchmarkHTTP_100Conn_5Delay(b *testing.B) {
	testRunBench(b, 100, 5)
}

func BenchmarkHTTP_100Conn_50Delay(b *testing.B) {
	testRunBench(b, 100, 50)
}

func testRunBench(b *testing.B, n int, blockMS int) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(blockMS))
	}))
	defer s.Close()

	c, _ := New(&Options{
		MaxConnsPerHost:    n,
		MaxConnWaitTimeout: time.Second * 10,
	})
	b.SetParallelism(20)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.Get(s.URL).doInternal()
			util.CheckError(err)
		}
	})
}
