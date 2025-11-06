package stock

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = errors.New("product not found")
	// ErrInsufficientStock is returned when there's not enough stock
	ErrInsufficientStock = errors.New("insufficient stock")
	// ErrInvalidQuantity is returned when quantity is invalid
	ErrInvalidQuantity = errors.New("invalid quantity")
)

// Store manages the stock inventory
type Store struct {
	mu       sync.RWMutex
	products map[string]*Product
}

// NewStore creates a new stock store
func NewStore() *Store {
	return &Store{
		products: make(map[string]*Product),
	}
}

// AddProduct adds a new product to the inventory
func (s *Store) AddProduct(product *Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product.UpdatedAt = time.Now()
	s.products[product.ID] = product
	return nil
}

// GetProduct retrieves a product by ID
func (s *Store) GetProduct(id string) (*Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.products[id]
	if !exists {
		return nil, ErrProductNotFound
	}
	return product, nil
}

// GetAllProducts returns all products
func (s *Store) GetAllProducts() []*Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]*Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}
	return products
}

// UpdateStock updates the stock quantity for a product
func (s *Store) UpdateStock(productID string, quantity int, operation string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[productID]
	if !exists {
		return ErrProductNotFound
	}

	switch operation {
	case "add":
		if quantity < 0 {
			return ErrInvalidQuantity
		}
		product.Quantity += quantity
	case "remove":
		if quantity < 0 {
			return ErrInvalidQuantity
		}
		if product.Quantity < quantity {
			return ErrInsufficientStock
		}
		product.Quantity -= quantity
	case "set":
		if quantity < 0 {
			return ErrInvalidQuantity
		}
		product.Quantity = quantity
	default:
		return errors.New("invalid operation: must be 'add', 'remove', or 'set'")
	}

	product.UpdatedAt = time.Now()
	return nil
}

// DeleteProduct removes a product from the inventory
func (s *Store) DeleteProduct(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[id]; !exists {
		return ErrProductNotFound
	}

	delete(s.products, id)
	return nil
}
