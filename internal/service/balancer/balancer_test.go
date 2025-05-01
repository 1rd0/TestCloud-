package balancer

import (
	"github.com/1rd0/TestCloud-/internal/service/backend"
	"net/url"
	"testing"
)

func TestRoundRobin_Next(t *testing.T) {
	u1, _ := url.Parse("http://localhost:8001")
	u2, _ := url.Parse("http://localhost:8002")
	b1 := backend.New(u1)
	b2 := backend.New(u2)

	rr := NewRR([]*backend.Backend{b1, b2})

	for i := 0; i < 10; i++ {
		b, err := rr.Next()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !b.IsAlive() {
			t.Errorf("backend not alive")
		}
	}
}

func BenchmarkRoundRobin_Next(b *testing.B) {
	u, _ := url.Parse("http://localhost:8001")
	back := backend.New(u)
	rr := NewRR([]*backend.Backend{back})

	for i := 0; i < b.N; i++ {
		_, _ = rr.Next()
	}
}
