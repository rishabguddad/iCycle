package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"macbook-marketplace/internal/catalog"
)

// handleHealth answers GET /api/health, so anyone can check the server is running.
func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleListCategories answers GET /api/categories with every category.
func (s *server) handleListCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, newCategoryResponses(s.store.Categories()))
}

// handleListCategoryProducts answers GET /api/categories/{id}/products.
// It accepts the same filters as GET /api/products.
func (s *server) handleListCategoryProducts(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "category")
	if !ok {
		return
	}
	if !s.store.CategoryExists(id) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}

	filter, err := parseFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	filter.CategoryID = id

	writeJSON(w, http.StatusOK, newProductResponses(s.store.Products(filter)))
}

// handleListProducts answers GET /api/products. Optional filters:
//
//	?search=battery  ?condition=used  ?inStock=true  ?categoryId=2
func (s *server) handleListProducts(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, newProductResponses(s.store.Products(filter)))
}

// handleGetProduct answers GET /api/products/{id} with one product.
func (s *server) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "product")
	if !ok {
		return
	}

	product, ok := s.store.Product(id)
	if !ok {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}

	writeJSON(w, http.StatusOK, newProductResponse(&product))
}

// parseFilter reads the optional filters from the query string
// (the part of the URL after "?").
func parseFilter(r *http.Request) (catalog.Filter, error) {
	query := r.URL.Query()
	filter := catalog.Filter{
		Search:    strings.TrimSpace(query.Get("search")),
		Condition: catalog.Condition(query.Get("condition")),
	}

	if filter.Condition != "" && !filter.Condition.Valid() {
		return catalog.Filter{}, errors.New("condition must be new, refurbished or used")
	}

	inStock := query.Get("inStock")
	if inStock != "" {
		value, err := strconv.ParseBool(inStock)
		if err != nil {
			return catalog.Filter{}, errors.New("inStock must be true or false")
		}
		filter.InStockOnly = value
	}

	categoryID := query.Get("categoryId")
	if categoryID != "" {
		value, err := strconv.Atoi(categoryID)
		if err != nil {
			return catalog.Filter{}, errors.New("categoryId must be a number")
		}
		filter.CategoryID = value
	}

	return filter, nil
}
