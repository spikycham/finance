package business

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spikycham/finance/network"
)

// MockRepository implements Repository interface for testing
type MockRepository struct {
	insertItemFunc   func(ctx context.Context, item *Item) error
	queryItemsFunc   func(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error)
	insertItemCalled int
	queryItemsCalled int
}

func (m *MockRepository) InsertItem(ctx context.Context, item *Item) error {
	m.insertItemCalled++
	if m.insertItemFunc != nil {
		return m.insertItemFunc(ctx, item)
	}
	return nil
}

func (m *MockRepository) QueryItemsByUserIDAndTime(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error) {
	m.queryItemsCalled++
	if m.queryItemsFunc != nil {
		return m.queryItemsFunc(ctx, userID, startTime)
	}
	return nil, nil
}

func TestNewService(t *testing.T) {
	mock := &MockRepository{}
	s := NewService(mock)

	if s == nil {
		t.Fatal("NewService() returned nil")
	}
	if s.r != mock {
		t.Error("NewService() did not set repository correctly")
	}
}

func TestService_CreateRecordItem(t *testing.T) {
	tests := []struct {
		name         string
		item         *Item
		mockFunc     func(ctx context.Context, item *Item) error
		wantErr      bool
		wantErrValue error
	}{
		{
			name: "successful creation",
			item: &Item{
				UserID:    uuid.New(),
				Type:      "food",
				Amount:    25.50,
				Note:      "lunch",
				CreatedAt: time.Now().UnixMilli(),
			},
			mockFunc:     nil,
			wantErr:      false,
			wantErrValue: nil,
		},
		{
			name: "repository error returns internal error",
			item: &Item{
				UserID:    uuid.New(),
				Type:      "food",
				Amount:    25.50,
				Note:      "lunch",
				CreatedAt: time.Now().UnixMilli(),
			},
			mockFunc: func(ctx context.Context, item *Item) error {
				return errors.New("database error")
			},
			wantErr:      true,
			wantErrValue: network.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRepository{
				insertItemFunc: tt.mockFunc,
			}
			s := NewService(mock)

			err := s.CreateRecordItem(context.Background(), tt.item)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecordItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !errors.Is(err, tt.wantErrValue) {
				t.Errorf("CreateRecordItem() error = %v, want %v", err, tt.wantErrValue)
			}

			if mock.insertItemCalled != 1 {
				t.Errorf("InsertItem called %d times, want 1", mock.insertItemCalled)
			}
		})
	}
}

func TestService_GetYearItems(t *testing.T) {
	userID := uuid.New()
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedItems := []Item{
		{
			ID:        1,
			UserID:    userID,
			Type:      "food",
			Amount:    25.50,
			Note:      "lunch",
			CreatedAt: baseTime.UnixMilli(),
		},
		{
			ID:        2,
			UserID:    userID,
			Type:      "transport",
			Amount:    10.00,
			Note:      "bus",
			CreatedAt: baseTime.AddDate(0, 1, 0).UnixMilli(),
		},
	}

	tests := []struct {
		name         string
		userID       uuid.UUID
		startTime    time.Time
		mockFunc     func(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error)
		wantItems    []Item
		wantErr      bool
		wantErrValue error
	}{
		{
			name:      "successful query",
			userID:    userID,
			startTime: baseTime,
			mockFunc: func(ctx context.Context, uid uuid.UUID, st time.Time) ([]Item, error) {
				return expectedItems, nil
			},
			wantItems:    expectedItems,
			wantErr:      false,
			wantErrValue: nil,
		},
		{
			name:      "empty result",
			userID:    userID,
			startTime: baseTime,
			mockFunc: func(ctx context.Context, uid uuid.UUID, st time.Time) ([]Item, error) {
				return []Item{}, nil
			},
			wantItems:    []Item{},
			wantErr:      false,
			wantErrValue: nil,
		},
		{
			name:      "repository error returns internal error",
			userID:    userID,
			startTime: baseTime,
			mockFunc: func(ctx context.Context, uid uuid.UUID, st time.Time) ([]Item, error) {
				return nil, errors.New("database error")
			},
			wantItems:    nil,
			wantErr:      true,
			wantErrValue: network.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRepository{
				queryItemsFunc: tt.mockFunc,
			}
			s := NewService(mock)

			got, err := s.GetYearItems(context.Background(), tt.userID, tt.startTime)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetYearItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !errors.Is(err, tt.wantErrValue) {
				t.Errorf("GetYearItems() error = %v, want %v", err, tt.wantErrValue)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.wantItems) {
					t.Errorf("GetYearItems() returned %d items, want %d", len(got), len(tt.wantItems))
					return
				}
				for i, item := range got {
					if item.ID != tt.wantItems[i].ID {
						t.Errorf("GetYearItems() item[%d].ID = %v, want %v", i, item.ID, tt.wantItems[i].ID)
					}
				}
			}

			if mock.queryItemsCalled != 1 {
				t.Errorf("QueryItemsByUserIDAndTime called %d times, want 1", mock.queryItemsCalled)
			}
		})
	}
}

func TestService_ContextPropagation(t *testing.T) {
	t.Run("context is passed to repository", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test-key"
		ctx := context.WithValue(context.Background(), testKey, "test-value")

		var capturedCtx context.Context
		mock := &MockRepository{
			insertItemFunc: func(ctx context.Context, item *Item) error {
				capturedCtx = ctx
				return nil
			},
		}
		s := NewService(mock)

		item := &Item{
			UserID:    uuid.New(),
			Type:      "food",
			Amount:    10.0,
			Note:      "test",
			CreatedAt: time.Now().UnixMilli(),
		}

		_ = s.CreateRecordItem(ctx, item)

		if capturedCtx.Value(testKey) != "test-value" {
			t.Error("Context was not properly propagated to repository")
		}
	})
}
