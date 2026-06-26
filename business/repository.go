package business

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func InsertItem() {}

// Get items in a full year.
func (r Repository) GetAllItemsByTime() {}
