package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/findardi/forflow-webhooks/server/internal/helper"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidURL           = errors.New("invalid url")
	ErrInvalidApplicationID = errors.New("invalid application_id")
	ErrInvalidEndpointID    = errors.New("invalid endpoint_id")
	ErrEndpointExists       = errors.New("endpoint url already registered for this application")
	ErrApplicationNotFound  = errors.New("application not found")
	ErrEndpointNotFound     = errors.New("endpoint not found")
)

type Service struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func validateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: scheme must be http or https", ErrInvalidURL)
	}

	if u.Hostname() == "" {
		return fmt.Errorf("%w: missing host", ErrInvalidURL)
	}

	return nil
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (sqlc.Endpoint, error) {
	if err := validateURL(request.URL); err != nil {
		return sqlc.Endpoint{}, err
	}

	appID, err := helper.ToUUID(request.ApplicationID)
	if err != nil {
		return sqlc.Endpoint{}, fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	result, err := s.q.CreateEndpoint(ctx, sqlc.CreateEndpointParams{
		ApplicationID: appID,
		Url:           request.URL,
	})

	if err != nil {
		// No row inserted: application missing or soft-deleted.
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.Endpoint{}, ErrApplicationNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return sqlc.Endpoint{}, ErrEndpointExists
		}

		return sqlc.Endpoint{}, fmt.Errorf("create endpoint: %w", err)
	}

	return result, nil
}

func (s *Service) ListByApplication(ctx context.Context, applicationID string) ([]sqlc.Endpoint, error) {
	appID, err := helper.ToUUID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	result, err := s.q.ListEndpointsByApplicationID(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("list endpoints: %w", err)
	}

	return result, nil
}

func (s *Service) Delete(ctx context.Context, applicationID, endpointID string) error {
	appID, err := helper.ToUUID(applicationID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	eID, err := helper.ToUUID(endpointID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidEndpointID, err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin delete: %w", err)
	}

	defer func() {
		_ = tx.Rollback(context.WithoutCancel(ctx))
	}()

	q := s.q.WithTx(tx)

	rows, err := q.DeleteEndpoint(ctx, sqlc.DeleteEndpointParams{ID: eID, ApplicationID: appID})
	if err != nil {
		return fmt.Errorf("delete endpoint: %w", err)
	}
	if rows == 0 {
		return ErrEndpointNotFound
	}

	if err := q.CancelStatusDeliveries(ctx, eID); err != nil {
		return fmt.Errorf("cancel deliveries status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit delete endpoint: %w", err)
	}

	return nil
}
