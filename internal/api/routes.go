// Package api turns HTTP requests into calls to the catalog,
// and the results into JSON replies.
package api

import (
	"net/http"

	"macbook-marketplace/internal/catalog"
)

// server holds what the handlers need. Each handler is a method on server,
// so it reaches the store through s.store instead of a global variable.
type server struct {
	store *catalog.Store
}

// NewRouter returns a router with every API route connected to its handler.
func NewRouter(store *catalog.Store) *http.ServeMux {
	s := &server{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/categories", s.handleListCategories)
	mux.HandleFunc("GET /api/categories/{id}/products", s.handleListCategoryProducts)
	mux.HandleFunc("GET /api/products", s.handleListProducts)
	mux.HandleFunc("GET /api/products/{id}", s.handleGetProduct)
	return mux
}
