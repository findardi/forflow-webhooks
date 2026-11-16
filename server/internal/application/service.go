package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/findardi/forflow-webhooks/server/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxNameLength = 100

var (
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidApplicationID = errors.New("invalid application_id")
	ErrApplicationNotFound  = errors.New("application not found")
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

func (s *Service) Create(ctx context.Context, request CreateRequest) (sqlc.Application, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return sqlc.Application{}, fmt.Errorf("%w: name is required", ErrInvalidName)
	}
	if utf8.RuneCountInString(name) > maxNameLength {
		return sqlc.Application{}, fmt.Errorf("%w: max %d characters", ErrInvalidName, maxNameLength)
	}

	result, err := s.q.CreateApplication(ctx, name)
	if err != nil {
		return sqlc.Application{}, fmt.Errorf("create application: %w", err)
	}

	return result, nil
}

func (s *Service) List(ctx context.Context) ([]sqlc.Application, error) {
	result, err := s.q.ListApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}

	return result, nil
}

func (s *Service) Delete(ctx context.Context, applicationID string) error {
	appID, err := helper.ToUUID(applicationID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin delete: %w", err)
	}

	defer func() {
		_ = tx.Rollback(context.WithoutCancel(ctx))
	}()

	q := s.q.WithTx(tx)

	rows, err := q.DeleteApplication(ctx, appID)
	if err != nil {
		return fmt.Errorf("delete application: %w", err)
	}
	if rows == 0 {
		return ErrApplicationNotFound
	}

	if err := q.DeleteEndpointsByApplication(ctx, appID); err != nil {
		return fmt.Errorf("delete endpoints: %w", err)
	}

	if err := q.CancelDeliveriesByApplication(ctx, appID); err != nil {
		return fmt.Errorf("cancel deliveries: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit delete: %w", err)
	}

	return nil
}
