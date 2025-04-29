package limiter

import (
	"sync"
	"time"
)

type bucket struct {
	tokens  int64 // текущее количество
	cap     int64 // ёмкость
	rate    int64 // токенов в секунду
	updated time.Time
	mu      sync.Mutex
}

func newBucket(cap, rate int64) *bucket { return nil }

func (b *bucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.updated)
	b.updated = now
	// пополняем
	add := (elapsed.Nanoseconds() * b.rate) / int64(time.Second)
	if add > 0 {
		b.tokens = min(b.cap, b.tokens+add)
	}
	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}
