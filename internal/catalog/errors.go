package catalog

import "errors"

// Named errors, so callers can tell failures apart with errors.Is
// (for example, the API answers 404 for one and 409 for another).
var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidQuantity = errors.New("quantity must be at least 1")
	ErrOutOfStock      = errors.New("not enough stock")
)
