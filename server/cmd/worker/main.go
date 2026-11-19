package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/findardi/forflow-webhooks/server/internal/app"
	"github.com/findardi/forflow-webhooks/server/internal/config"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		log.Error("load config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.RunWorker(ctx, cfg, log); err != nil {
		log.Error("run worker", "err", err)
		stop()
		os.Exit(1)
	}
}
