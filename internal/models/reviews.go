package models

type (
	Review struct {
		ID          int    `json:"id" db:"id"`
		Stars       string `json:"stars" db:"stars"`
		Body        string `json:"body" db:"body"`
		ProductID   string `json:"product_id" db:"product_id"`
		ProductName string `json:"product_name" db:"product_name"`
	}

	ReviewsFilter struct {
		Query *string `json:"query"`
	}
)
