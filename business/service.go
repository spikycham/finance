package business

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/spikycham/finance/logger"
	"github.com/spikycham/finance/network"
)

type Repository interface {
	InsertItem(ctx context.Context, item *Item) error
	QueryItemsByUserIDAndTime(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time) ([]Item, error)
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) CreateRecordItem(ctx context.Context, item *Item) error {
	if err := s.r.InsertItem(ctx, item); err != nil {
		logger.Error("InsertItem failed", err)
		return network.ErrInternal
	}

	return nil
}

func (s *Service) GetYearItems(ctx context.Context, userID uuid.UUID, year int) ([]Item, error) {
	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

	items, err := s.r.QueryItemsByUserIDAndTime(ctx, userID, startTime, endTime)
	if err != nil {
		logger.Error("QueryItemsByUserIDAndTime failed", err)
		return nil, network.ErrInternal
	}

	return items, nil
}
