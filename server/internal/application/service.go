package application

import (
	"context"
	"fmt"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
	}
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (sqlc.Application, error) {
	result, err := s.q.CreateApplication(ctx, request.Name)
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
