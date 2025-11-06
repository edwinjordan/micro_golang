package stock

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Handler provides HTTP handlers for the stock service
type Handler struct {
	store *Store
}

// NewHandler creates a new handler
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// ServeHTTP routes requests to appropriate handlers
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set common headers
	w.Header().Set("Content-Type", "application/json")

	// Route based on path and method
	path := strings.TrimPrefix(r.URL.Path, "/api/stock")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			h.handleGetAllProducts(w, r)
		case http.MethodPost:
			h.handleAddProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle /api/stock/{id} and /api/stock/{id}/update
	if strings.HasSuffix(path, "/update") {
		if r.Method == http.MethodPost {
			h.handleUpdateStock(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle /api/stock/{id}
	switch r.Method {
	case http.MethodGet:
		h.handleGetProduct(w, r)
	case http.MethodDelete:
		h.handleDeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	products := h.store.GetAllProducts()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"count":    len(products),
	})
}

func (h *Handler) handleAddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if product.ID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	if err := h.store.AddProduct(&product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	log.Printf("Added product: %s", product.ID)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/stock/")
	productID := strings.TrimSuffix(path, "/")

	product, err := h.store.GetProduct(productID)
	if err != nil {
		if err == ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *Handler) handleUpdateStock(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/stock/")
	productID := strings.TrimSuffix(path, "/update")

	var update StockUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	update.ProductID = productID

	if err := h.store.UpdateStock(update.ProductID, update.Quantity, update.Operation); err != nil {
		if err == ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else if err == ErrInsufficientStock {
			http.Error(w, "Insufficient stock", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	product, _ := h.store.GetProduct(productID)
	json.NewEncoder(w).Encode(product)
	log.Printf("Updated stock for product %s: %s %d", productID, update.Operation, update.Quantity)
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/stock/")
	productID := strings.TrimSuffix(path, "/")

	if err := h.store.DeleteProduct(productID); err != nil {
		if err == ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Deleted product: %s", productID)
}
