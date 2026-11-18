// Package api turns HTTP requests into calls to the catalog and order
// packages, and their results into JSON replies.
package api

import (
	"net/http"

	"macbook-marketplace/internal/catalog"
	"macbook-marketplace/internal/order"
)

// server holds what the handlers need. Each handler is a method on server,
// so it reaches the store and orders through s instead of global variables.
type server struct {
	store  *catalog.Store
	orders *order.Service
}

// NewRouter returns a router with every API route connected to its handler.
func NewRouter(store *catalog.Store, orders *order.Service) *http.ServeMux {
	s := &server{store: store, orders: orders}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)

	mux.HandleFunc("GET /api/categories", s.handleListCategories)
	mux.HandleFunc("GET /api/categories/{id}/products", s.handleListCategoryProducts)

	mux.HandleFunc("GET /api/products", s.handleListProducts)
	mux.HandleFunc("GET /api/products/{id}", s.handleGetProduct)

	mux.HandleFunc("POST /api/orders", s.handlePlaceOrder)
	mux.HandleFunc("GET /api/orders", s.handleListOrders)
	mux.HandleFunc("GET /api/orders/{id}", s.handleGetOrder)
	mux.HandleFunc("POST /api/orders/{id}/cancel", s.handleCancelOrder)
	return mux
}
