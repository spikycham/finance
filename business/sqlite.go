package business

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func Connect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize the record items table and index on user ID in database.
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'others',
			amount REAL NOT NULL,
			note TEXT,
			created_at INTEGER NOT NULL
		)`,
	); err != nil {
		return nil, err
	}
	if _, err := db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_items_user_id
		ON items(user_id)`,
	); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *SQLiteRepository) InsertItem(ctx context.Context, item *Item) error {
	if _, err := r.db.ExecContext(
		ctx,
		`INSERT INTO items (user_id, type, amount, note, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		item.UserID,
		item.Type,
		item.Amount,
		item.Note,
		item.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

// Get items in a full year.
func (r *SQLiteRepository) QueryItemsByUserIDAndTime(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time) ([]Item, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, amount, note, created_at
		FROM items
		WHERE user_id = ? AND created_at >= ? AND created_at < ?
		ORDER BY created_at`,
		userID,
		startTime.Unix(),
		endTime.Unix(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Type,
			&item.Amount,
			&item.Note,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
