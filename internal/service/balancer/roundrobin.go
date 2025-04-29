package balancer

import (
	"errors"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	back []*url.URL
	idx  uint64
}

func NewRR(raw []string) (*RoundRobin, error) {
	urls := make([]*url.URL, 0, len(raw))
	for _, s := range raw {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	if len(urls) == 0 {
		return nil, errors.New("no backends supplied")
	}
	return &RoundRobin{back: urls}, nil
}

func (r *RoundRobin) Next() (*url.URL, error) {
	i := atomic.AddUint64(&r.idx, 1)
	return r.back[int(i-1)%len(r.back)], nil
}
