package business

import (
	"database/sql"
)

func Connect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL DEFAULT 'others',
			amount REAL NOT NULL,
			note TEXT,
			created_at INTEGER NOT NULL
		)`,
	); err != nil {
		return nil, err
	}

	return db, nil
}
