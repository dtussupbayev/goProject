package models

type (
	Product struct {
		ID          int    `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
		Price       int    `json:"price" db:"price"`
		CategoryID  string `json:"category_id" db:"category_id"`
	}

	ProductsFilter struct {
		Query *string `json:"query"`
	}
)
