package main

import (
	"context"
	"github.com/1rd0/TestCloud-/internal/server"
	"log"
	"os/signal"

	"syscall"
)

func main() {
	cfgPath := "config/config.yaml"

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, cfgPath); err != nil {
		log.Fatal(err)
	}
}
