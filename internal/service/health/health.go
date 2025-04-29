package health

import (
	"context"
	"github.com/1rd0/TestCloud-/internal/service/metrics"
	"net/http"
	"net/url"
	"time"

	"github.com/1rd0/TestCloud-/internal/service/backend"
	"go.uber.org/zap"
)

func Start(ctx context.Context, backs []*backend.Backend, interval, timeout time.Duration, log *zap.Logger) {
	client := &http.Client{Timeout: timeout}
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, b := range backs {
					go func(b *backend.Backend) {
						alive := ping(client, b.URL)
						if alive != b.IsAlive() {
							log.Info("backend state changed",
								zap.String("url", b.URL.String()),
								zap.Bool("alive", alive))
						}
						b.SetAlive(alive)
						val := float64(0)
						if alive {
							val = 1
						}

						metrics.BackendUp.WithLabelValues(b.URL.String()).Set(val)
					}(b)
				}
			}
		}
	}()
}

func ping(c *http.Client, u *url.URL) bool {
	resp, err := c.Get(u.String() + "/health")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
