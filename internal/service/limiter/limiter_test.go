package limiter

import (
	"testing"
)

func TestBucket_TryTake(t *testing.T) {
	b := NewBucket(2, 1)
	if err := b.TryTake(); err != nil {
		t.Fatal("expected no error on first take")
	}
	if err := b.TryTake(); err != nil {
		t.Fatal("expected no error on second take")
	}
	if err := b.TryTake(); err == nil {
		t.Fatal("expected error on third take due to rate limit")
	}
}

func BenchmarkBucket_TryTake(b *testing.B) {
	bucket := NewBucket(100, 100)
	for i := 0; i < b.N; i++ {
		_ = bucket.TryTake()
	}
}
