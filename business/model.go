package business

// Types: food, supply, transport, treatment, study, job, wear, workout, entertainment, house, income, others.
type Item struct {
	ID        int64
	Type      string
	Amount    float64
	Note      string
	CreatedAt int64
}
