package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/findardi/forflow-webhooks/server/db"
	"github.com/findardi/forflow-webhooks/server/internal/config"
	"github.com/findardi/forflow-webhooks/server/internal/endpoint"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ShutdownTimeout = 10 * time.Second

func RunAPI(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	log.Info("database connection good")
	defer pool.Close()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           NewRouter(pool),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("http server listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}

		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	log.Info("shutting down http server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

func NewRouter(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	EndpointRouter(mux, pool)

	return mux
}

func EndpointRouter(mux *http.ServeMux, pool *pgxpool.Pool) {
	h := endpoint.NewHandler(pool)

	mux.HandleFunc("POST /endpoints", h.CreateEndpoint)
}
