package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	url   *url.URL
	conns int64
	proxy *httputil.ReverseProxy
}

func newBackend(u *url.URL) *Backend {
	rp := httputil.NewSingleHostReverseProxy(u)
	return &Backend{url: u, proxy: rp}
}

func (b *Backend) URL() *url.URL    { return b.url }
func (b *Backend) ConnCount() int64 { return atomic.LoadInt64(&b.conns) }

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&b.conns, 1)
	defer atomic.AddInt64(&b.conns, -1)
	b.proxy.ServeHTTP(w, r)
}
