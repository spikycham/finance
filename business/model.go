package business

import (
	"github.com/google/uuid"
)

// Types: food, supply, transport, treatment, study, job, wear, workout, entertainment, house, income, others.
type Item struct {
	ID        int64     `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Note      string    `json:"note"`
	CreatedAt int64     `json:"created_at"`
}
