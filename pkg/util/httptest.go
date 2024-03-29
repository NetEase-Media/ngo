package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func PerformRequest(r http.Handler, method, path string, headers ...Header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func RunRequest(b *testing.B, r http.Handler, method, path string, headers ...Header) {
	// create fake request
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

type Header struct {
	Key   string
	Value string
}
