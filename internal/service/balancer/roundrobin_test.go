package integration

import (
	"github.com/1rd0/TestCloud-/internal/balancer"
	"testing"
)

func TestRR(t *testing.T) {
	rr, _ := balancer.NewRR([]string{"http://a", "http://b", "http://c"})
	want := []string{"http://a", "http://b", "http://c", "http://a"}

	for i, exp := range want {
		if got := rr.Next().String(); got != exp {
			t.Fatalf("%d: want %s, got %s", i, exp, got)
		}
	}
}
