package main

import (
	"context"
	"github.com/1rd0/TestCloud-/internal/server"
	"os/signal"

	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		panic(err)
	}
}
