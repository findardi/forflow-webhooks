package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/findardi/forflow-webhooks/server/internal/helper"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidURL           = errors.New("invalid url")
	ErrInvalidApplicationID = errors.New("invalid application_id")
	ErrEndpointExists       = errors.New("endpoint url already registered for this application")
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
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
