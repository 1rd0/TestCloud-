package balancer

import (
	"errors"

	"github.com/1rd0/TestCloud-/internal/service/backend"
)

type LeastConn struct{ back []*backend.Backend }

func NewLC(back []*backend.Backend) *LeastConn { return &LeastConn{back: back} }

func (l *LeastConn) Next() (*backend.Backend, error) {
	var best *backend.Backend
	for _, b := range l.back {
		if !b.IsAlive() {
			continue
		}
		if best == nil || b.ConnCount() < best.ConnCount() {
			best = b
		}
	}
	if best == nil {
		return nil, errors.New("no alive backends")
	}
	return best, nil
}
