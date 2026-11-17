package api

import (
	"net/http"
	"strconv"

	"macbook-marketplace/internal/catalog"
)

// handleHealth answers GET /api/health, so anyone can check the server is running.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleListCategories answers GET /api/categories with every category.
func handleListCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, catalog.Categories())
}

// handleListCategoryProducts answers GET /api/categories/{id}/products
// with the products in one category.
func handleListCategoryProducts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "category id must be a number")
		return
	}

	if !catalog.CategoryExists(id) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}

	writeJSON(w, http.StatusOK, catalog.ProductsInCategory(id))
}

// handleListProducts answers GET /api/products with every product.
func handleListProducts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, catalog.Products())
}

// handleGetProduct answers GET /api/products/{id} with one product.
func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "product id must be a number")
		return
	}

	product := catalog.FindProduct(id)
	if product == nil {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}

	writeJSON(w, http.StatusOK, product)
}
