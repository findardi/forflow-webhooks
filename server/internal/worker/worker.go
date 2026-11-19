package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
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
	jobs := make(chan sqlc.GetClaimRow)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(jobs)
		return w.poll(ctx, jobs)
	})

	for range 5 {
		g.Go(func() error {
			for job := range jobs {
				workCtx := context.WithoutCancel(ctx)
				w.svc.Record(workCtx, job, w.svc.Send(workCtx, job))
			}
			return nil
		})
	}

	return g.Wait()
}

func (w *Worker) poll(ctx context.Context, jobs chan<- sqlc.GetClaimRow) error {
	for {
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

		jobs <- job
	}
}
