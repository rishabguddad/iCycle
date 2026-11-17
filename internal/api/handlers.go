package api

import (
	"net/http"
	"strconv"
)

// handleHealth answers GET /api/health, so anyone can check the server is running.
func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleListCategories answers GET /api/categories with every category.
func (s *server) handleListCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, newCategoryResponses(s.store.Categories()))
}

// handleListCategoryProducts answers GET /api/categories/{id}/products
// with the products in one category.
func (s *server) handleListCategoryProducts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "category id must be a number")
		return
	}

	if !s.store.CategoryExists(id) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}

	writeJSON(w, http.StatusOK, newProductResponses(s.store.ProductsInCategory(id)))
}

// handleListProducts answers GET /api/products with every product.
func (s *server) handleListProducts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, newProductResponses(s.store.Products()))
}

// handleGetProduct answers GET /api/products/{id} with one product.
func (s *server) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "product id must be a number")
		return
	}

	product, ok := s.store.Product(id)
	if !ok {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}

	writeJSON(w, http.StatusOK, newProductResponse(&product))
}
