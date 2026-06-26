package business

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/spikycham/finance/network"
)

type Repository interface {
	InsertItem(ctx context.Context, item *Item) error
	QueryItemsByUserIDAndTime(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error)
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) CreateRecordItem(ctx context.Context, item *Item) error {
	if err := s.r.InsertItem(ctx, item); err != nil {
		// TODO: add a logger
		return network.ErrInternal
	}

	return nil
}

func (s *Service) GetYearItems(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error) {
	items, err := s.r.QueryItemsByUserIDAndTime(ctx, userID, startTime)
	if err != nil {
		// TODO: add a logger
		return nil, network.ErrInternal
	}

	return items, nil
}

// PERF: update and delete items from records.
