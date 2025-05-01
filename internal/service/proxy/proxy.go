package proxy

import (
	"github.com/1rd0/TestCloud-/pkg/metrics"
	"net/http"

	"github.com/1rd0/TestCloud-/internal/service/backend"
)

type Handler struct {
	pick func() (*backend.Backend, error)
}

func New(pick func() (*backend.Backend, error)) *Handler {
	return &Handler{
		pick: pick,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/metrics" || r.URL.Path == "/health" || r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	metrics.TotalRequests.Inc()

	b, err := h.pick()
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	b.ServeHTTP(w, r)
}
