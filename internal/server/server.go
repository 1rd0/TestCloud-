package server

import (
	"context"
	"github.com/1rd0/TestCloud-/config"
	"github.com/1rd0/TestCloud-/internal/service/backend"
	"github.com/1rd0/TestCloud-/internal/service/balancer"
	"github.com/1rd0/TestCloud-/internal/service/health"
	"github.com/1rd0/TestCloud-/internal/service/limiter"
	"github.com/1rd0/TestCloud-/internal/service/proxy"
	"github.com/1rd0/TestCloud-/pkg/gp"
	"github.com/1rd0/TestCloud-/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Run(ctx context.Context, path string) error {

	cfg, err := config.New(path)
	if err != nil {
		return err
	}
	log, err := logger.New()
	if err != nil {
		return err
	}
	backs := make([]*backend.Backend, 0, len(cfg.LB.Backends))
	for _, raw := range cfg.LB.Backends {
		if !strings.Contains(raw, "://") {
			raw = "http://" + raw
		}
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		backs = append(backs, backend.New(u))
	}
	pool, err := gp.NewPoolConn(ctx, cfg.DB.URL())
	if err != nil {
		return err
	}
	rateLimiterm, err := limiter.NewLimiter(ctx, pool)
	// choose algorithm

	bal := balancer.NewRR(backs)

	// health-check
	health.Start(ctx, backs, 5*time.Second, 2*time.Second, log)

	// HTTP-handler
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	apiHandler := proxy.New(bal.Next)
	mux.Handle("/", rateLimiterm.Middleware(apiHandler))
	srv := &http.Server{
		Addr:    cfg.Listen, // ":8040"
		Handler: mux,
	}

	log.Info("LB listening", zap.String("addr", cfg.Listen))
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Info("Shutting down server...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("Server shutdown failed", zap.Error(err))
		} else {
			log.Info("Server shutdown completed")
		}

	}()

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
