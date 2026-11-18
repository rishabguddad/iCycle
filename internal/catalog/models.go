package catalog

import (
	"fmt"
	"strings"
)

// Condition is the state a product is sold in. It is a named type instead
// of a plain string, so the allowed values live in one place and can be
// checked with Valid.
type Condition string

const (
	ConditionNew         Condition = "new"
	ConditionRefurbished Condition = "refurbished"
	ConditionUsed        Condition = "used"
)

// Valid reports whether c is one of the allowed conditions.
func (c Condition) Valid() bool {
	switch c {
	case ConditionNew, ConditionRefurbished, ConditionUsed:
		return true
	}
	return false
}

// Category groups products, e.g. MacBooks or Accessories.
type Category struct {
	ID   int
	Name string
}

// ProductDetails is the part of a product an admin is allowed to edit.
// Stock is not here on purpose: it changes only through sales, cancellations
// and restocks, never by typing a new number.
type ProductDetails struct {
	Name             string
	Description      string
	CategoryID       int
	Condition        Condition
	Price            int64 // in paise
	Specs            map[string]string
	CompatibleModels []string
}

// validate checks the details. categoryIDs holds every category that exists.
func (d ProductDetails) validate(categoryIDs map[int]bool) error {
	switch {
	case strings.TrimSpace(d.Name) == "":
		return fmt.Errorf("%w: name is required", ErrInvalidProduct)
	case !categoryIDs[d.CategoryID]:
		return fmt.Errorf("%w: category %d does not exist", ErrInvalidProduct, d.CategoryID)
	case !d.Condition.Valid():
		return fmt.Errorf("%w: condition %q must be new, refurbished or used", ErrInvalidProduct, d.Condition)
	case d.Price <= 0:
		return fmt.Errorf("%w: price must be more than 0", ErrInvalidProduct)
	}
	return nil
}

// Product is anything the store sells: a MacBook, a part or an accessory.
//
// It embeds ProductDetails. The fields of ProductDetails become fields of
// Product, so p.Name and p.Price work directly, and an update can replace
// every editable field at once with p.ProductDetails = newDetails.
type Product struct {
	ID int
	ProductDetails

	// stock is private: other packages can read it with Stock(), but can
	// only change it through Sell and Restock, which enforce the rules
	// (never negative, never more sold than available).
	stock int
}
