package limiter

import (
	"errors"
	"sync"
	"time"
)

type Bucket struct {
	capacity     int
	tokens       float64
	refillPerSec float64
	last         time.Time
	mu           sync.Mutex
}

func NewBucket(capacity int, refillPerSec float64) *Bucket {
	return &Bucket{
		capacity:     capacity,
		tokens:       float64(capacity),
		refillPerSec: refillPerSec,
		last:         time.Now(),
	}
}

func (b *Bucket) TryTake() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	delta := now.Sub(b.last).Seconds() * b.refillPerSec
	b.tokens += delta
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	b.last = now

	if b.tokens >= 1 {
		b.tokens -= 1
		return nil
	}
	return errors.New("rate limit exceeded")
}
