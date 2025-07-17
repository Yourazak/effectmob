package subscriptions

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, sub Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *Service) List(ctx context.Context) ([]Subscription, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int) (*Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, sub Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) SumPrice(ctx context.Context, userIDStr, serviceName string, start, end time.Time) (int, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 0, err
	}
	return s.repo.SumPrice(ctx, userID, serviceName, start, end)
}
