package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidURL           = errors.New("invalid url")
	ErrInvalidApplicationID = errors.New("invalid application_id")
	ErrAlreadyUsed          = errors.New("already exists")
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
	}
}

func toUUID(u string) (pgtype.UUID, error) {
	var id pgtype.UUID

	if u == "" {
		return id, nil
	}

	if err := id.Scan(u); err != nil {
		return id, fmt.Errorf("parse to uuid: %w", err)
	}

	return id, nil
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

func (s *Service) CreateEndpoint(ctx context.Context, request CreateEndpoint) (sqlc.Endpoint, error) {
	if err := validateURL(request.Url); err != nil {
		return sqlc.Endpoint{}, err
	}

	appID, err := toUUID(request.ApplicationID)
	if err != nil {
		return sqlc.Endpoint{}, fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	result, err := s.q.CreateEndpoint(ctx, sqlc.CreateEndpointParams{
		ApplicationID: appID,
		Url:           request.Url,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.Endpoint{}, ErrAlreadyUsed
		}

		return sqlc.Endpoint{}, fmt.Errorf("create endpoint: %w", err)
	}

	return result, nil
}
