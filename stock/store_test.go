package stock_test

import (
	"testing"
	"time"

	"github.com/edwinjordan/micro_golang/stock"
)

func TestNewStore(t *testing.T) {
	store := stock.NewStore()
	if store == nil {
		t.Fatal("Expected store to be created, got nil")
	}
}

func TestAddProduct(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:          "TEST001",
		Name:        "Test Product",
		Description: "A test product",
		Quantity:    10,
		Price:       99.99,
	}

	err := store.AddProduct(product)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify product was added
	retrieved, err := store.GetProduct("TEST001")
	if err != nil {
		t.Fatalf("Expected to retrieve product, got error: %v", err)
	}

	if retrieved.ID != product.ID {
		t.Errorf("Expected product ID %s, got %s", product.ID, retrieved.ID)
	}

	if retrieved.Name != product.Name {
		t.Errorf("Expected product name %s, got %s", product.Name, retrieved.Name)
	}

	if retrieved.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	store := stock.NewStore()

	_, err := store.GetProduct("NONEXISTENT")
	if err != stock.ErrProductNotFound {
		t.Errorf("Expected ErrProductNotFound, got %v", err)
	}
}

func TestGetAllProducts(t *testing.T) {
	store := stock.NewStore()

	// Add multiple products
	products := []*stock.Product{
		{ID: "P1", Name: "Product 1", Quantity: 5, Price: 10.00},
		{ID: "P2", Name: "Product 2", Quantity: 10, Price: 20.00},
		{ID: "P3", Name: "Product 3", Quantity: 15, Price: 30.00},
	}

	for _, p := range products {
		store.AddProduct(p)
	}

	allProducts := store.GetAllProducts()
	if len(allProducts) != 3 {
		t.Errorf("Expected 3 products, got %d", len(allProducts))
	}
}

func TestUpdateStock_Add(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Add stock
	err := store.UpdateStock("TEST001", 5, "add")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updated, _ := store.GetProduct("TEST001")
	if updated.Quantity != 15 {
		t.Errorf("Expected quantity 15, got %d", updated.Quantity)
	}
}

func TestUpdateStock_Remove(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Remove stock
	err := store.UpdateStock("TEST001", 3, "remove")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updated, _ := store.GetProduct("TEST001")
	if updated.Quantity != 7 {
		t.Errorf("Expected quantity 7, got %d", updated.Quantity)
	}
}

func TestUpdateStock_InsufficientStock(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 5,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Try to remove more than available
	err := store.UpdateStock("TEST001", 10, "remove")
	if err != stock.ErrInsufficientStock {
		t.Errorf("Expected ErrInsufficientStock, got %v", err)
	}
}

func TestUpdateStock_Set(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Set stock to specific value
	err := store.UpdateStock("TEST001", 25, "set")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updated, _ := store.GetProduct("TEST001")
	if updated.Quantity != 25 {
		t.Errorf("Expected quantity 25, got %d", updated.Quantity)
	}
}

func TestUpdateStock_InvalidQuantity(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Try negative quantity
	err := store.UpdateStock("TEST001", -5, "add")
	if err != stock.ErrInvalidQuantity {
		t.Errorf("Expected ErrInvalidQuantity, got %v", err)
	}
}

func TestUpdateStock_InvalidOperation(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Try invalid operation
	err := store.UpdateStock("TEST001", 5, "invalid")
	if err == nil {
		t.Error("Expected error for invalid operation, got nil")
	}
}

func TestUpdateStock_ProductNotFound(t *testing.T) {
	store := stock.NewStore()

	err := store.UpdateStock("NONEXISTENT", 5, "add")
	if err != stock.ErrProductNotFound {
		t.Errorf("Expected ErrProductNotFound, got %v", err)
	}
}

func TestDeleteProduct(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Delete the product
	err := store.DeleteProduct("TEST001")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify product is deleted
	_, err = store.GetProduct("TEST001")
	if err != stock.ErrProductNotFound {
		t.Errorf("Expected ErrProductNotFound after delete, got %v", err)
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	store := stock.NewStore()

	err := store.DeleteProduct("NONEXISTENT")
	if err != stock.ErrProductNotFound {
		t.Errorf("Expected ErrProductNotFound, got %v", err)
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := stock.NewStore()

	product := &stock.Product{
		ID:       "TEST001",
		Name:     "Test Product",
		Quantity: 100,
		Price:    99.99,
	}
	store.AddProduct(product)

	// Test concurrent reads and writes
	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				store.GetProduct("TEST001")
			}
			done <- true
		}()
	}

	// Concurrent writes
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				store.UpdateStock("TEST001", 1, "add")
				time.Sleep(time.Millisecond)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 15; i++ {
		<-done
	}

	// Verify final state
	final, _ := store.GetProduct("TEST001")
	if final.Quantity != 200 {
		t.Errorf("Expected final quantity 200, got %d", final.Quantity)
	}
}
