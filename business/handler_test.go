package business

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spikycham/finance/network"
)

// MockService implements a mock service for handler testing
type MockService struct {
	createRecordItemFunc func(ctx context.Context, item *Item) error
	getYearItemsFunc     func(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error)
	createCalled         int
	getCalled            int
}

func (m *MockService) CreateRecordItem(ctx context.Context, item *Item) error {
	m.createCalled++
	if m.createRecordItemFunc != nil {
		return m.createRecordItemFunc(ctx, item)
	}
	return nil
}

func (m *MockService) GetYearItems(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error) {
	m.getCalled++
	if m.getYearItemsFunc != nil {
		return m.getYearItemsFunc(ctx, userID, startTime)
	}
	return nil, nil
}

// Ensure MockService satisfies the Service interface methods used by Handler
// We need to adapt the Handler to accept an interface instead of concrete *Service

func TestNewHandler(t *testing.T) {
	mockRepo := &MockRepository{}
	s := NewService(mockRepo)
	h := NewHandler(s)

	if h == nil {
		t.Fatal("NewHandler() returned nil")
	}
	if h.s != s {
		t.Error("NewHandler() did not set service correctly")
	}
}

func TestHandler_CreateItem(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		body           string
		wantCode       int
		wantMessage    string
		wantError      string
		mockFunc       func(ctx context.Context, item *Item) error
	}{
		{
			name:        "successful creation",
			body:        `{"user_id":"` + userID.String() + `","type":"food","amount":25.50,"note":"lunch"}`,
			wantCode:    http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:    nil,
		},
		{
			name:      "invalid JSON body",
			body:      `{invalid json`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrInvalidJSON.Error(),
			mockFunc:  nil,
		},
		{
			name:      "missing user_id",
			body:      `{"type":"food","amount":25.50,"note":"lunch"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "missing type",
			body:      `{"user_id":"` + userID.String() + `","amount":25.50,"note":"lunch"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "invalid type value",
			body:      `{"user_id":"` + userID.String() + `","type":"invalid","amount":25.50,"note":"lunch"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "missing amount",
			body:      `{"user_id":"` + userID.String() + `","type":"food","note":"lunch"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "zero amount",
			body:      `{"user_id":"` + userID.String() + `","type":"food","amount":0,"note":"lunch"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "negative amount",
			body:      `{"user_id":"` + userID.String() + `","type":"food","amount":-10,"note":"lunch"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "note too long",
			body:      `{"user_id":"` + userID.String() + `","type":"food","amount":25.50,"note":"` + strings.Repeat("a", 101) + `"}`,
			wantCode:  http.StatusBadRequest,
			wantError: network.ErrMissingFields.Error(),
			mockFunc:  nil,
		},
		{
			name:      "all valid types - supply",
			body:      `{"user_id":"` + userID.String() + `","type":"supply","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - transport",
			body:      `{"user_id":"` + userID.String() + `","type":"transport","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - treatment",
			body:      `{"user_id":"` + userID.String() + `","type":"treatment","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - study",
			body:      `{"user_id":"` + userID.String() + `","type":"study","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - job",
			body:      `{"user_id":"` + userID.String() + `","type":"job","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - wear",
			body:      `{"user_id":"` + userID.String() + `","type":"wear","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - workout",
			body:      `{"user_id":"` + userID.String() + `","type":"workout","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - entertainment",
			body:      `{"user_id":"` + userID.String() + `","type":"entertainment","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - house",
			body:      `{"user_id":"` + userID.String() + `","type":"house","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - income",
			body:      `{"user_id":"` + userID.String() + `","type":"income","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "all valid types - others",
			body:      `{"user_id":"` + userID.String() + `","type":"others","amount":10,"note":"test"}`,
			wantCode:  http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:  nil,
		},
		{
			name:      "service error returns 500",
			body:      `{"user_id":"` + userID.String() + `","type":"food","amount":25.50,"note":"lunch"}`,
			wantCode:  http.StatusInternalServerError,
			wantError: "internal error",
			mockFunc: func(ctx context.Context, item *Item) error {
				return errors.New("database error")
			},
		},
		{
			name:        "empty note is valid",
			body:        `{"user_id":"` + userID.String() + `","type":"food","amount":25.50,"note":""}`,
			wantCode:    http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:    nil,
		},
		{
			name:        "max length note is valid",
			body:        `{"user_id":"` + userID.String() + `","type":"food","amount":25.50,"note":"` + strings.Repeat("a", 100) + `"}`,
			wantCode:    http.StatusCreated,
			wantMessage: "success to create a new financial record",
			mockFunc:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			s := NewService(mockRepo)

			// Override the mock behavior for this test
			if tt.mockFunc != nil {
				// We need to use a custom repository that calls the mock function
				mockRepo.insertItemFunc = tt.mockFunc
			}

			h := NewHandler(s)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/create", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateItem(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("CreateItem() status = %v, want %v", w.Code, tt.wantCode)
			}

			var resp network.StandardResponse[any]
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.wantMessage != "" {
				if resp.Message == nil || *resp.Message != tt.wantMessage {
					t.Errorf("CreateItem() message = %v, want %v", resp.Message, tt.wantMessage)
				}
			}

			if tt.wantError != "" {
				if resp.Error == nil || *resp.Error != tt.wantError {
					t.Errorf("CreateItem() error = %v, want %v", resp.Error, tt.wantError)
				}
			}
		})
	}
}

func TestHandler_GetItems(t *testing.T) {
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
	}

	tests := []struct {
		name        string
		queryParams string
		wantCode    int
		wantError   string
		wantItems   []Item
		mockFunc    func(ctx context.Context, userID uuid.UUID, startTime time.Time) ([]Item, error)
	}{
		{
			name:        "successful query",
			queryParams: "?user_id=" + userID.String() + "&start_time=" + strconv.FormatInt(baseTime.UnixMilli()-1, 10),
			wantCode:    http.StatusOK,
			wantItems:   expectedItems,
			mockFunc: func(ctx context.Context, uid uuid.UUID, st time.Time) ([]Item, error) {
				return expectedItems, nil
			},
		},
		{
			name:        "missing user_id",
			queryParams: "?start_time=1735689600000",
			wantCode:    http.StatusBadRequest,
			wantError:   network.ErrMissingFields.Error(),
			mockFunc:    nil,
		},
		{
			name:        "invalid user_id format",
			queryParams: "?user_id=invalid&start_time=1735689600000",
			wantCode:    http.StatusBadRequest,
			wantError:   network.ErrMissingFields.Error(),
			mockFunc:    nil,
		},
		{
			name:        "missing start_time",
			queryParams: "?user_id=" + userID.String(),
			wantCode:    http.StatusBadRequest,
			wantError:   network.ErrMissingFields.Error(),
			mockFunc:    nil,
		},
		{
			name:        "invalid start_time format",
			queryParams: "?user_id=" + userID.String() + "&start_time=invalid",
			wantCode:    http.StatusBadRequest,
			wantError:   network.ErrMissingFields.Error(),
			mockFunc:    nil,
		},
		{
			name:        "service error returns 500",
			queryParams: "?user_id=" + userID.String() + "&start_time=1735689600000",
			wantCode:    http.StatusInternalServerError,
			wantError:   "internal error",
			mockFunc: func(ctx context.Context, uid uuid.UUID, st time.Time) ([]Item, error) {
				return nil, errors.New("database error")
			},
		},
		{
			name:        "empty result",
			queryParams: "?user_id=" + userID.String() + "&start_time=1735689600000",
			wantCode:    http.StatusOK,
			wantItems:   []Item{},
			mockFunc: func(ctx context.Context, uid uuid.UUID, st time.Time) ([]Item, error) {
				return []Item{}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				queryItemsFunc: tt.mockFunc,
			}
			s := NewService(mockRepo)
			h := NewHandler(s)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/items"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			h.GetItems(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("GetItems() status = %v, want %v", w.Code, tt.wantCode)
			}

			var resp network.StandardResponse[GetItemsResponse]
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.wantError != "" {
				if resp.Error == nil || *resp.Error != tt.wantError {
					t.Errorf("GetItems() error = %v, want %v", resp.Error, tt.wantError)
				}
			}

			if tt.wantItems != nil && resp.Data != nil {
				if len(resp.Data.Items) != len(tt.wantItems) {
					t.Errorf("GetItems() returned %d items, want %d", len(resp.Data.Items), len(tt.wantItems))
				}
			}
		})
	}
}
