package catalog

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

// clone returns a full copy of the product. Specs and CompatibleModels are
// copied too, so changing the copy can never change the original.
func (p *Product) clone() Product {
	c := *p
	c.Specs = maps.Clone(p.Specs)
	c.CompatibleModels = slices.Clone(p.CompatibleModels)
	return c
}

// Stock returns how many units are available.
func (p *Product) Stock() int {
	return p.stock
}

// InStock reports whether at least one unit is available.
func (p *Product) InStock() bool {
	return p.stock > 0
}

// Sell reduces the product's stock by quantity.
// It returns an error, and changes nothing, if the sale isn't possible.
func (p *Product) Sell(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be at least 1")
	}
	if quantity > p.stock {
		return fmt.Errorf("only %d in stock", p.stock)
	}
	p.stock = p.stock - quantity
	return nil
}

// Restock adds quantity units, e.g. when a delivery arrives
// or an order is cancelled.
func (p *Product) Restock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be at least 1")
	}
	p.stock = p.stock + quantity
	return nil
}
