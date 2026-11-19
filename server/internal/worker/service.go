package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxDrainBytes = 64 << 10

type Service struct {
	q          *sqlc.Queries
	httpClient *http.Client
	Log        *slog.Logger
}

func NewService(pool *pgxpool.Pool, log *slog.Logger) *Service {
	return &Service{
		q: sqlc.New(pool),
		httpClient: &http.Client{
			Timeout:       30 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		Log: log,
	}
}

func (s *Service) Claim(ctx context.Context) (sqlc.GetClaimRow, bool, error) {
	job, err := s.q.GetClaim(ctx)

	if errors.Is(err, pgx.ErrNoRows) {
		s.Log.Info("no job pending")
		return sqlc.GetClaimRow{}, false, nil
	}

	if err != nil {
		return sqlc.GetClaimRow{}, false, fmt.Errorf("claim delivery: %w", err)
	}

	return job, true, nil
}

func (s *Service) Send(ctx context.Context, job sqlc.GetClaimRow) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, job.Url, bytes.NewReader(job.Payload))
	if err != nil {
		s.Log.Warn("build request", "delivery_id", job.DeliveryID.String(), "error", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", job.MessageID.String())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.Log.Warn("send request", "delivery_id", job.DeliveryID.String(), "error", err)
		return false
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrainBytes))

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (s *Service) Record(ctx context.Context, job sqlc.GetClaimRow, ok bool) {
	mark := s.q.MarkFailed
	if ok {
		mark = s.q.MarkSuccess
	}

	rows, err := mark(ctx, job.DeliveryID)
	if err != nil {
		s.Log.Error("record delivery", "delivery_id", job.DeliveryID.String(), "error", err)
		return
	}

	if rows == 0 {
		s.Log.Info("delivery no longer pending", "delivery_id", job.DeliveryID.String())
		return
	}

	s.Log.Info("delivery recorded", "delivery_id", job.DeliveryID.String(), "ok", ok)
}
