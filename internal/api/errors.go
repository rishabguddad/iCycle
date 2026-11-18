package api

import (
	"errors"
	"net/http"

	"macbook-marketplace/internal/catalog"
	"macbook-marketplace/internal/order"
)

// writeServiceError picks the HTTP status for an error returned by the
// catalog or order package. errors.Is also finds errors wrapped with %w.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	// 400: the request itself is wrong, and would be wrong at any time.
	case errors.Is(err, catalog.ErrInvalidProduct),
		errors.Is(err, catalog.ErrInvalidQuantity),
		errors.Is(err, order.ErrInvalid):
		writeError(w, http.StatusBadRequest, err.Error())

	// 404: the thing asked for doesn't exist.
	case errors.Is(err, catalog.ErrProductNotFound),
		errors.Is(err, order.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	// 409: the request is valid, but clashes with the current state.
	case errors.Is(err, catalog.ErrOutOfStock),
		errors.Is(err, order.ErrAlreadyCancelled):
		writeError(w, http.StatusConflict, err.Error())

	default:
		writeError(w, http.StatusInternalServerError, "something went wrong")
	}
}
