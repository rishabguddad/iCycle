package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// handleHealth answers GET /api/health, so anyone can check the server is running.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleListCategories answers GET /api/categories with every category.
func handleListCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, categories)
}

// handleListCategoryProducts answers GET /api/categories/{id}/products
// with the products in one category.
func handleListCategoryProducts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "category id must be a number")
		return
	}

	if !categoryExists(id) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}

	writeJSON(w, http.StatusOK, productsInCategory(id))
}

// handleListProducts answers GET /api/products with every product.
func handleListProducts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, products)
}

// handleGetProduct answers GET /api/products/{id} with one product.
func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		writeError(w, http.StatusBadRequest, "product id must be a number")
		return
	}

	product := findProduct(id)
	if product == nil {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}

	writeJSON(w, http.StatusOK, product)
}

// writeError sends an error message as JSON, e.g. {"error":"product not found"}.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// writeJSON sends data to the client as JSON, with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
