package message

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
	ErrInvalidEndpointID    = errors.New("invalid endpoint_id")
	ErrInvalidPayload       = errors.New("payload is required")
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

func (s *Service) CreateMessage(ctx context.Context, request CreateMessage) (sqlc.Message, error) {
	if len(request.Payload) == 0 || string(request.Payload) == "null" {
		return sqlc.Message{}, ErrInvalidPayload
	}

	appID, err := helper.ToUUID(request.ApplicationID)
	if err != nil {
		return sqlc.Message{}, fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return sqlc.Message{}, fmt.Errorf("begin create message: %w", err)
	}

	defer func() {
		_ = tx.Rollback(context.WithoutCancel(ctx))
	}()

	q := s.q.WithTx(tx)

	message, err := q.CreateMessages(ctx, sqlc.CreateMessagesParams{
		ApplicationID: appID,
		Payload:       request.Payload,
	})

	if err != nil {
		// No row inserted: application missing or soft-deleted.
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.Message{}, ErrApplicationNotFound
		}

		return sqlc.Message{}, fmt.Errorf("create message: %w", err)
	}

	if _, err := q.CreateDeliveriesForMessage(ctx, message.ID); err != nil {
		return sqlc.Message{}, fmt.Errorf("create deliveries for message: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.Message{}, fmt.Errorf("commit create message: %w", err)
	}

	return message, nil
}

func (s *Service) ListByEndpoint(ctx context.Context, applicationID, endpointID string) ([]sqlc.Message, error) {
	appID, err := helper.ToUUID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidApplicationID, err)
	}

	eID, err := helper.ToUUID(endpointID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidEndpointID, err)
	}

	result, err := s.q.ListMessagesByEndpoint(ctx, sqlc.ListMessagesByEndpointParams{
		EndpointID:    eID,
		ApplicationID: appID,
	})
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}

	return result, nil
}
