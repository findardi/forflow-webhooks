package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
)

const maxNameLength = 100

var ErrInvalidName = errors.New("invalid name")

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
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
