package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

const (
	pollInterval   = time.Second
	workers        = 50
	maxPerEndpoint = 2
)

type Worker struct {
	svc *Service
}

func NewWorker(pool *pgxpool.Pool, log *slog.Logger) *Worker {
	return &Worker{
		svc: NewService(pool, log),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	n, err := w.svc.q.ResetProcessing(ctx)
	if err != nil {
		return fmt.Errorf("reset processing: %w", err)
	}

	if n > 0 {
		w.svc.Log.Warn("orpahned deliveries reset to pending", "count", n)
	}

	jobs := make(chan sqlc.GetClaimRow)
	sem := make(chan struct{}, workers)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(jobs)
		return w.poll(ctx, jobs, sem)
	})

	for range workers {
		g.Go(func() error {
			for job := range jobs {
				workCtx := context.WithoutCancel(ctx)
				w.svc.Record(workCtx, job, w.svc.Send(workCtx, job))
				<-sem
			}
			return nil
		})
	}

	return g.Wait()
}

func (w *Worker) poll(ctx context.Context, jobs chan<- sqlc.GetClaimRow, sem chan struct{}) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case sem <- struct{}{}:
		}

		job, ok, err := w.svc.Claim(ctx, maxPerEndpoint)
		if err != nil {
			<-sem
			if ctx.Err() != nil {
				return nil
			}

			return err
		}

		if !ok {
			<-sem
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
