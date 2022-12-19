package models

//type Category struct {
//	ID   int    `json:"id" db:"id"`
//	Name string `json:"name" db:"name"`
//}

//type Product struct {
//	ID          int    `json:"id" db:"id"`
//	Name        string `json:"name" db:"name"`
//	Description string `json:"description" db:"description"`
//	Price       int    `json:"price" db:"price"`
//	CategoryID  string `json:"category_id" db:"category_id"`
//}

type Order struct {
	ID     int    `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
}
