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
// Admin routes, and routes that show customers' details, need adminPassword.
func NewRouter(store *catalog.Store, orders *order.Service, adminPassword string) *http.ServeMux {
	s := &server{store: store, orders: orders}

	// admin wraps a handler so it needs the admin login.
	admin := func(handler http.HandlerFunc) http.Handler {
		return AdminOnly(adminPassword, handler)
	}

	mux := http.NewServeMux()

	// Public: anyone can browse and place an order.
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/categories", s.handleListCategories)
	mux.HandleFunc("GET /api/categories/{id}/products", s.handleListCategoryProducts)
	mux.HandleFunc("GET /api/products", s.handleListProducts)
	mux.HandleFunc("GET /api/products/{id}", s.handleGetProduct)
	mux.HandleFunc("POST /api/orders", s.handlePlaceOrder)

	// Admin only: manage the catalog.
	mux.Handle("POST /api/products", admin(s.handleCreateProduct))
	mux.Handle("PUT /api/products/{id}", admin(s.handleUpdateProduct))
	mux.Handle("DELETE /api/products/{id}", admin(s.handleDeleteProduct))
	mux.Handle("POST /api/products/{id}/restock", admin(s.handleRestockProduct))

	// Admin only: orders hold customers' names, phones and addresses.
	mux.Handle("GET /api/orders", admin(s.handleListOrders))
	mux.Handle("GET /api/orders/{id}", admin(s.handleGetOrder))
	mux.Handle("POST /api/orders/{id}/cancel", admin(s.handleCancelOrder))
	return mux
}
