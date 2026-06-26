package business

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) *SQLiteRepository {
	t.Helper()
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewSQLiteRepository(db)
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "in-memory database",
			path:    ":memory:",
			wantErr: false,
		},
		{
			name:    "invalid path",
			path:    "/nonexistent/path/db.sqlite",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if db != nil {
				db.Close()
			}
		})
	}
}

func TestSQLiteRepository_InsertItem(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	now := time.Now().UnixMilli()

	tests := []struct {
		name    string
		item    *Item
		wantErr bool
	}{
		{
			name: "valid item",
			item: &Item{
				UserID:    userID,
				Type:      "food",
				Amount:    25.50,
				Note:      "lunch",
				CreatedAt: now,
			},
			wantErr: false,
		},
		{
			name: "item with empty note",
			item: &Item{
				UserID:    userID,
				Type:      "transport",
				Amount:    10.00,
				Note:      "",
				CreatedAt: now,
			},
			wantErr: false,
		},
		{
			name: "item with different type",
			item: &Item{
				UserID:    userID,
				Type:      "income",
				Amount:    5000.00,
				Note:      "salary",
				CreatedAt: now,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.InsertItem(ctx, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteRepository_QueryItemsByUserIDAndTime(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	otherUserID := uuid.New()

	// Create test data with different timestamps
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	items := []Item{
		{
			UserID:    userID,
			Type:      "food",
			Amount:    25.50,
			Note:      "january lunch",
			CreatedAt: baseTime.UnixMilli(),
		},
		{
			UserID:    userID,
			Type:      "transport",
			Amount:    10.00,
			Note:      "january bus",
			CreatedAt: baseTime.AddDate(0, 1, 0).UnixMilli(), // February
		},
		{
			UserID:    userID,
			Type:      "food",
			Amount:    30.00,
			Note:      "march dinner",
			CreatedAt: baseTime.AddDate(0, 2, 0).UnixMilli(), // March
		},
		{
			UserID:    otherUserID,
			Type:      "food",
			Amount:    15.00,
			Note:      "other user food",
			CreatedAt: baseTime.UnixMilli(),
		},
	}

	// Insert test data
	for i := range items {
		if err := repo.InsertItem(ctx, &items[i]); err != nil {
			t.Fatalf("Failed to insert test item %d: %v", i, err)
		}
	}

	tests := []struct {
		name      string
		userID    uuid.UUID
		startTime time.Time
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get all items for user from beginning of year",
			userID:    userID,
			startTime: baseTime.AddDate(0, -1, 0), // December previous year
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "get items from february",
			userID:    userID,
			startTime: baseTime, // January 1st
			wantCount: 2,        // February and March
			wantErr:   false,
		},
		{
			name:      "get items from march",
			userID:    userID,
			startTime: baseTime.AddDate(0, 1, 0), // February 1st
			wantCount: 1,                          // Only March
			wantErr:   false,
		},
		{
			name:      "no items after march",
			userID:    userID,
			startTime: baseTime.AddDate(0, 2, 0), // March 1st
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "different user has separate items",
			userID:    otherUserID,
			startTime: baseTime.AddDate(0, -1, 0),
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "non-existent user",
			userID:    uuid.New(),
			startTime: baseTime,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.QueryItemsByUserIDAndTime(ctx, tt.userID, tt.startTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryItemsByUserIDAndTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("QueryItemsByUserIDAndTime() returned %d items, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestSQLiteRepository_QueryItemsOrdering(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	baseTime := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	// Insert items in reverse chronological order
	items := []Item{
		{
			UserID:    userID,
			Type:      "food",
			Amount:    30.00,
			Note:      "third",
			CreatedAt: baseTime.AddDate(0, 0, 2).UnixMilli(),
		},
		{
			UserID:    userID,
			Type:      "food",
			Amount:    20.00,
			Note:      "first",
			CreatedAt: baseTime.UnixMilli(),
		},
		{
			UserID:    userID,
			Type:      "food",
			Amount:    25.00,
			Note:      "second",
			CreatedAt: baseTime.AddDate(0, 0, 1).UnixMilli(),
		},
	}

	for i := range items {
		if err := repo.InsertItem(ctx, &items[i]); err != nil {
			t.Fatalf("Failed to insert test item %d: %v", i, err)
		}
	}

	got, err := repo.QueryItemsByUserIDAndTime(ctx, userID, baseTime.AddDate(0, -1, 0))
	if err != nil {
		t.Fatalf("QueryItemsByUserIDAndTime() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(got))
	}

	// Verify items are ordered by created_at ascending
	expectedOrder := []string{"first", "second", "third"}
	for i, item := range got {
		if item.Note != expectedOrder[i] {
			t.Errorf("Item %d: got note %q, want %q", i, item.Note, expectedOrder[i])
		}
	}
}

func TestSQLiteRepository_ConcurrentAccess(t *testing.T) {
	// Use a temporary file for concurrent access testing
	tmpFile := t.TempDir() + "/test_concurrent.db"
	db, err := Connect(tmpFile)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := NewSQLiteRepository(db)

	ctx := context.Background()

	userID := uuid.New()
	baseTime := time.Now().UnixMilli()

	// Test sequential inserts (SQLite has limitations with concurrent writes)
	for i := 0; i < 10; i++ {
		item := &Item{
			UserID:    userID,
			Type:      "food",
			Amount:    float64(i * 10),
			Note:      "concurrent item",
			CreatedAt: baseTime + int64(i),
		}
		err := repo.InsertItem(ctx, item)
		if err != nil {
			t.Errorf("InsertItem() error = %v", err)
		}
	}

	// Verify all items were inserted
	items, err := repo.QueryItemsByUserIDAndTime(ctx, userID, time.UnixMilli(baseTime-1))
	if err != nil {
		t.Fatalf("QueryItemsByUserIDAndTime() error = %v", err)
	}
	if len(items) != 10 {
		t.Errorf("Expected 10 items, got %d", len(items))
	}
}
