package stock

import "time"

// Product represents a product in the inventory
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockUpdate represents a stock quantity update
type StockUpdate struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Operation string `json:"operation"` // "add" or "remove"
}
