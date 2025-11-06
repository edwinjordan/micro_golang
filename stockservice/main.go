package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edwinjordan/micro_golang/stock"
)

const (
	defaultPort = "8080"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create stock store and handler
	store := stock.NewStore()
	handler := stock.NewHandler(store)

	// Initialize with some sample data
	initializeSampleData(store)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.Handle("/api/stock/", handler)
	mux.Handle("/api/stock", handler)

	// Add a health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Stock management service starting on port %s", port)
		log.Println("Available endpoints:")
		log.Println("  GET    /api/stock       - List all products")
		log.Println("  POST   /api/stock       - Add new product")
		log.Println("  GET    /api/stock/{id}  - Get product by ID")
		log.Println("  DELETE /api/stock/{id}  - Delete product")
		log.Println("  POST   /api/stock/{id}/update - Update stock quantity")
		log.Println("  GET    /health          - Health check")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down stock service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Stock service stopped")
}

func initializeSampleData(store *stock.Store) {
	sampleProducts := []*stock.Product{
		{
			ID:          "LAPTOP001",
			Name:        "Dell Laptop XPS 15",
			Description: "High-performance laptop for professionals",
			Quantity:    50,
			Price:       1299.99,
		},
		{
			ID:          "MOUSE001",
			Name:        "Wireless Mouse",
			Description: "Ergonomic wireless mouse",
			Quantity:    200,
			Price:       29.99,
		},
		{
			ID:          "KEYBOARD001",
			Name:        "Mechanical Keyboard",
			Description: "RGB mechanical gaming keyboard",
			Quantity:    75,
			Price:       89.99,
		},
	}

	for _, product := range sampleProducts {
		if err := store.AddProduct(product); err != nil {
			log.Printf("Warning: Failed to add sample product %s: %v", product.ID, err)
		}
	}

	log.Printf("Initialized with %d sample products", len(sampleProducts))
}
