package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/findardi/forflow-webhooks/server/internal/netguard"
	"github.com/findardi/forflow-webhooks/server/internal/signing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var retrySchedule = []time.Duration{
	5 * time.Second,
	5 * time.Minute,
	30 * time.Minute,
	2 * time.Hour,
	5 * time.Hour,
	10 * time.Hour,
}

const maxDrainBytes = 64 << 10

type Service struct {
	q          *sqlc.Queries
	httpClient *http.Client
	Log        *slog.Logger
	Pool       *pgxpool.Pool
	now        func() time.Time
	jitter     func(time.Duration) time.Duration
}

func NewService(pool *pgxpool.Pool, log *slog.Logger, development string) *Service {
	var transport http.RoundTripper = netguard.NewTransport()
	if development == "dev" {
		transport = http.DefaultTransport
	}

	return &Service{
		q:    sqlc.New(pool),
		Pool: pool,
		httpClient: &http.Client{
			Timeout:       15 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
			Transport:     transport,
		},
		Log:    log,
		now:    time.Now,
		jitter: tenPercentJitter,
	}
}

func tenPercentJitter(d time.Duration) time.Duration {
	if d < 10 {
		return d
	}
	return d + rand.N(d/10)
}

func (s *Service) Claim(ctx context.Context, maxPerEndpoint int32) (sqlc.GetClaimRow, bool, error) {
	job, err := s.q.GetClaim(ctx, maxPerEndpoint)

	if errors.Is(err, pgx.ErrNoRows) {
		s.Log.Debug("no job pending")
		return sqlc.GetClaimRow{}, false, nil
	}

	if err != nil {
		return sqlc.GetClaimRow{}, false, fmt.Errorf("claim delivery: %w", err)
	}

	return job, true, nil
}

type outcome int

const (
	outcomeSent outcome = iota
	outcomeEndpointFailed
	outcomeDiscarded
)

type sendResult struct {
	outcome        outcome
	responseStatus *int32
	errorKind      string
}

func errorKind(err error) string {
	var dnsErr *net.DNSError
	var netErr net.Error
	switch {
	case errors.Is(err, netguard.ErrBlocked):
		return "blocked"
	case errors.As(err, &dnsErr):
		return "dns"
	case errors.As(err, &netErr) && netErr.Timeout():
		return "timeout"
	case errors.Is(err, syscall.ECONNREFUSED):
		return "conn_refused"
	}
	return "other"
}

func (s *Service) Send(ctx context.Context, job sqlc.GetClaimRow) sendResult {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, job.Url, bytes.NewReader(job.Payload))
	if err != nil {
		s.Log.Warn("build request", "delivery_id", job.DeliveryID.String())
		return sendResult{
			outcome:   outcomeDiscarded,
			errorKind: "invalid_request",
		}
	}

	ts := s.now()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", job.MessageID.String())
	req.Header.Set("webhook-timestamp", signing.FormatTimestamp(ts))

	req.Header.Set("webhook-signature", signing.Sign(job.Secret, job.MessageID.String(), ts, job.Payload))
	resp, err := s.httpClient.Do(req)
	if err != nil {
		kind := errorKind(err)
		s.Log.Warn("send request", "delivery_id", job.DeliveryID.String(), "error_kind", kind)
		return sendResult{
			outcome:   outcomeEndpointFailed,
			errorKind: kind,
		}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrainBytes))

	code := int32(resp.StatusCode)
	result := sendResult{
		outcome:        outcomeEndpointFailed,
		responseStatus: &code,
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.outcome = outcomeSent
	}

	return result
}

func nextRetry(attemptsBefore int32) (wait time.Duration, ok bool) {
	if attemptsBefore < 0 || int(attemptsBefore) >= len(retrySchedule) {
		return 0, false
	}

	return retrySchedule[attemptsBefore], true
}

func (s *Service) Record(ctx context.Context, job sqlc.GetClaimRow, res sendResult) {
	var code pgtype.Int4
	if res.responseStatus != nil {
		code = pgtype.Int4{
			Int32: *res.responseStatus,
			Valid: true,
		}
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		s.Log.Error("begin record", "delivery_id", job.DeliveryID.String(), "error", err)
		return
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()
	q := s.q.WithTx(tx)

	var (
		rows   int64
		status string
	)

	switch res.outcome {
	case outcomeSent:
		status = "succeeded"
		rows, err = q.MarkSuccess(ctx, job.DeliveryID)
	case outcomeDiscarded:
		status = "failed"
		rows, err = q.MarkFailed(ctx, job.DeliveryID)
	case outcomeEndpointFailed:
		if wait, ok := nextRetry(job.AttemptCount); ok {
			status = "retry"
			next := s.now().Add(s.jitter(wait)).UTC()
			rows, err = q.RetryAttempt(ctx, sqlc.RetryAttemptParams{
				ID:            job.DeliveryID,
				NextAttemptAt: pgtype.Timestamptz{Time: next, Valid: true},
			})
		} else {
			status = "failed"
			rows, err = q.MarkFailed(ctx, job.DeliveryID)
		}
	}

	if err == nil && rows > 0 {
		var kind pgtype.Text
		if res.errorKind != "" {
			kind = pgtype.Text{
				String: res.errorKind,
				Valid:  true,
			}
		}

		if err = q.CreateAttempt(ctx, sqlc.CreateAttemptParams{
			DeliveryID:     job.DeliveryID,
			ResponseStatus: code,
			ErrorKind:      kind,
		}); err != nil {
			err = fmt.Errorf("insert attempt: %w", err)
		}
	}

	if err == nil && rows > 0 {
		if err = tx.Commit(ctx); err != nil {
			err = fmt.Errorf("commit record: %w", err)
		}
	}

	attrs := []any{
		"delivery_id", job.DeliveryID.String(),
		"attempt", job.AttemptCount + 1,
		"outcome", status,
		"error_kind", res.errorKind,
	}

	if res.responseStatus != nil {
		attrs = append(attrs, "response_status", *res.responseStatus)
	}

	switch {
	case err != nil:
		s.Log.Error("record delivery outcome", append(attrs, "error", err)...)
	case rows == 0:
		s.Log.Warn("lease lost before recording outcome", attrs...)
	default:
		s.Log.Info("delivery attempt recorded", attrs...)
	}
}
