package api

import (
	"net/http"
	"strings"

	"macbook-marketplace/internal/catalog"
)

// productRequest is the body of PUT /api/products/{id}: the editable details.
// There is no stock field, so sending "stock" is rejected as an unknown field.
type productRequest struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	CategoryID       int               `json:"categoryId"`
	Condition        string            `json:"condition"`
	Price            int64             `json:"price"` // in paise
	Image            string            `json:"image"`
	Specs            map[string]string `json:"specs"`
	CompatibleModels []string          `json:"compatibleModels"`
}

// createProductRequest is the body of POST /api/products. It embeds
// productRequest, so it has all those fields, plus a starting stock.
type createProductRequest struct {
	productRequest
	Stock int `json:"stock"`
}

func (p productRequest) toDetails() catalog.ProductDetails {
	return catalog.ProductDetails{
		Name:             p.Name,
		Description:      p.Description,
		CategoryID:       p.CategoryID,
		Condition:        catalog.Condition(p.Condition),
		Price:            p.Price,
		Image:            strings.TrimSpace(p.Image),
		Specs:            p.Specs,
		CompatibleModels: p.CompatibleModels,
	}
}

// handleCreateProduct answers POST /api/products.
func (s *server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var request createProductRequest
	err := decodeJSON(w, r, &request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	product, err := s.store.AddProduct(request.toDetails(), request.Stock)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newProductResponse(&product))
}

// handleUpdateProduct answers PUT /api/products/{id}.
func (s *server) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "product")
	if !ok {
		return
	}

	var request productRequest
	err := decodeJSON(w, r, &request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	product, err := s.store.UpdateProduct(id, request.toDetails())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newProductResponse(&product))
}

// handleDeleteProduct answers DELETE /api/products/{id}.
// On success it replies 204 No Content: done, and nothing to send back.
func (s *server) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "product")
	if !ok {
		return
	}

	err := s.store.DeleteProduct(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleRestockProduct answers POST /api/products/{id}/restock
// with a body like {"quantity": 10}.
func (s *server) handleRestockProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "product")
	if !ok {
		return
	}

	var request struct {
		Quantity int `json:"quantity"`
	}
	err := decodeJSON(w, r, &request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	product, err := s.store.Restock(id, request.Quantity)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newProductResponse(&product))
}
