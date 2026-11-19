package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

const pollInterval = time.Second

type Worker struct {
	svc *Service
}

func NewWorker(pool *pgxpool.Pool, log *slog.Logger) *Worker {
	return &Worker{
		svc: NewService(pool, log),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	for range 5 {
		wg.Go(func() error {
			return w.Loop(ctx)
		})
	}

	return wg.Wait()
}

func (w *Worker) Loop(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return nil
		}

		job, ok, err := w.svc.Claim(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return err
		}

		if !ok {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(pollInterval):
			}
			continue
		}

		// A claimed job finishes even during shutdown.
		workCtx := context.WithoutCancel(ctx)
		sent := w.svc.Send(workCtx, job)
		w.svc.Record(workCtx, job, sent)
	}
}
