package delivery

import (
	"context"
	"errors"
	"fmt"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/findardi/forflow-webhooks/server/internal/helper"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidApplicationID = errors.New("invalid application_id")
	ErrInvalidDeliveryID    = errors.New("invalid delivery_id")
	ErrDeliveryNotFound     = errors.New("delivery not found")
)

type Service struct {
	q *sqlc.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{q: sqlc.New(pool)}
}

func (s *Service) Get(ctx context.Context, applicationID, deliveryID string) (sqlc.GetDeliveryByApplicationRow, error) {
	appID, err := helper.ToUUID(applicationID)
	if err != nil {
		return sqlc.GetDeliveryByApplicationRow{}, fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	dID, err := helper.ToUUID(deliveryID)
	if err != nil {
		return sqlc.GetDeliveryByApplicationRow{}, fmt.Errorf("%w: %w", ErrInvalidDeliveryID, err)
	}

	d, err := s.q.GetDeliveryByApplication(ctx, sqlc.GetDeliveryByApplicationParams{
		ID:            dID,
		ApplicationID: appID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.GetDeliveryByApplicationRow{}, ErrDeliveryNotFound
	}
	if err != nil {
		return sqlc.GetDeliveryByApplicationRow{}, fmt.Errorf("get delivery: %w", err)
	}

	return d, nil
}

func (s *Service) ListAttempts(ctx context.Context, applicationID, deliveryID string) ([]sqlc.Attempt, error) {
	d, err := s.Get(ctx, applicationID, deliveryID)
	if err != nil {
		return nil, err
	}

	attempts, err := s.q.ListAttemptsByDelivery(ctx, d.ID)
	if err != nil {
		return nil, fmt.Errorf("list attempts: %w", err)
	}

	return attempts, nil
}
