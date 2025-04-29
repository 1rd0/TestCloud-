package balancer

import (
	"errors"
	"github.com/1rd0/TestCloud-/internal/service/backend"

	"sync/atomic"
)

type RoundRobin struct {
	back []*backend.Backend
	idx  uint64
}

func NewRR(back []*backend.Backend) *RoundRobin { return &RoundRobin{back: back} }

func (r *RoundRobin) Next() (*backend.Backend, error) {
	n := len(r.back)
	for i := 0; i < n; i++ {
		b := r.back[int(atomic.AddUint64(&r.idx, 1)-1)%n]
		if b.IsAlive() {
			return b, nil
		}
	}
	return nil, errors.New("no alive backends")
}
