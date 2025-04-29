package server

import (
	"context"
	"github.com/1rd0/TestCloud-/config"
	"github.com/1rd0/TestCloud-/internal/service/balancer"
	"github.com/1rd0/TestCloud-/internal/service/health"

	"github.com/1rd0/TestCloud-/internal/service/proxy"

	"go.uber.org/zap"

	"time"

	"net/http"
)

func Run(ctx context.Context) error {
	cfg, err := config.New("config/config.yaml")
	if err != nil {

	}

	log, err := zap.NewProduction()
	if err != nil {

	}
	rr, err := balancer.NewRR(cfg.LB.Backends) // round-robin
	if err != nil {

	}
	health.Start(ctx, cfg.LB.Backends, 5*time.Second, log)
	handler := proxy.New(rr)
	srv := &http.Server{
		Addr:    cfg.Listen,
		Handler: handler,
	}

	log.Info("LB listening", zap.String("addr", cfg.Listen))
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}
