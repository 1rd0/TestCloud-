package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL   *url.URL
	proxy *httputil.ReverseProxy

	alive uint32 // 0 | 1
	conns int64
}

func New(u *url.URL) *Backend {
	return &Backend{
		URL:   u,
		proxy: httputil.NewSingleHostReverseProxy(u),
		alive: 1,
	}
}

func (b *Backend) IsAlive() bool { return atomic.LoadUint32(&b.alive) == 1 }
func (b *Backend) SetAlive(ok bool) {
	val := uint32(0)
	if ok {
		val = 1
	}
	atomic.StoreUint32(&b.alive, val)
}
func (b *Backend) ConnCount() int64 { return atomic.LoadInt64(&b.conns) }

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	atomic.AddInt64(&b.conns, 1)
	defer atomic.AddInt64(&b.conns, -1)
	b.proxy.ServeHTTP(w, r)
}
