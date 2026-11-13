package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/findardi/forflow-webhooks/server/db"
	"github.com/findardi/forflow-webhooks/server/internal/config"
)

const ShutdownTimeout = 10 * time.Second

func RunAPI(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	log.Info("database connection good")
	defer pool.Close()

	return nil
}
