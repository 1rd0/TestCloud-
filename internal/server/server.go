package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/1rd0/TestCloud-/config"
	"github.com/1rd0/TestCloud-/internal/service/backend"
	"github.com/1rd0/TestCloud-/internal/service/balancer"
	"github.com/1rd0/TestCloud-/internal/service/health"
	"github.com/1rd0/TestCloud-/internal/service/proxy"

	"go.uber.org/zap"
)

func Run(ctx context.Context, path string) error {

	cfg, err := config.New(path)
	if err != nil {
		log.Fatal(err)
	}
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
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

	// choose algorithm
	var bal balancer.Balancer
	bal = balancer.NewRR(backs)

	// health-check
	health.Start(ctx, backs, 5*time.Second, 2*time.Second, log)

	// HTTP-handler
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/", proxy.New(bal.Next))
	mux.HandleFunc("/favicon.ico", http.NotFound)
	srv := &http.Server{
		Addr:    cfg.Listen, // ":8080"
		Handler: mux,
	}

	log.Info("LB listening", zap.String("addr", cfg.Listen))
	go func() { <-ctx.Done(); _ = srv.Shutdown(context.Background()) }()

	return srv.ListenAndServe()
}
