package proxy

import (
	"github.com/1rd0/TestCloud-/internal/balancer"
	"net/http"
	"net/http/httputil"
)

type Handler struct {
	bal balancer.Balancer
}

func New(bal balancer.Balancer) *Handler {
	return &Handler{bal: bal}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := h.bal.Next()
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	// Пересобираем request → целится в backend.
	r.URL.Scheme = u.Scheme
	r.URL.Host = u.Host
	r.Host = u.Host
	// Довершаем проксирование стандартным reverse-proxy’ем.
	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
