// Package api turns HTTP requests into calls to the catalog,
// and the results into JSON replies.
package api

import "net/http"

// NewRouter returns a router with every API route connected to its handler.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/categories", handleListCategories)
	mux.HandleFunc("GET /api/categories/{id}/products", handleListCategoryProducts)
	mux.HandleFunc("GET /api/products", handleListProducts)
	mux.HandleFunc("GET /api/products/{id}", handleGetProduct)
	return mux
}
