package proxy

/*type Backend struct {
	url   *url.URL
	alive uint32        // 0 | 1
	conns int64         // active connections
	rp    *httputil.ReverseProxy
}

func newBackend(u *url.URL) *Backend { … }

func (b *Backend) IsAlive() bool        { return atomic.LoadUint32(&b.alive) == 1 }
func (b *Backend) SetAlive(ok bool)     { val := uint32(0); if ok { val = 1 }; atomic.StoreUint32(&b.alive, val) }
func (b *Backend) URL() *url.URL        { return b.url }
func (b *Backend) ConnCount() int64     { return atomic.LoadInt64(&b.conns) }

func (b *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&b.conns, 1)
	defer atomic.AddInt64(&b.conns, -1)
	b.rp.ServeHTTP(w, r)
}
func New(backURLs []string, log *zap.Logger) ([]*backend.Backend, error) {
	var backs []*backend.Backend
	for _, raw := range backURLs {
		u, err := url.Parse(raw)
		if err != nil { return nil, err }
		rp := httputil.NewSingleHostReverseProxy(u)

		// error handler — делаем один ретрай
		rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			if retry, _ := r.Context().Value("retry").(bool); retry {
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}
			ctx := context.WithValue(r.Context(), "retry", true)
			rp.ServeHTTP(w, r.WithContext(ctx))
		}
		backs = append(backs, &backend.Backend{url: u, rp: rp, alive: 1})
	}
	return backs, nil
}*/
