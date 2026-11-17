package catalog

import (
	"errors"
	"fmt"
)

// Sell reduces the product's stock by quantity.
// It returns an error, and changes nothing, if the sale isn't possible.
func (p *Product) Sell(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be at least 1")
	}
	if quantity > p.Stock {
		return fmt.Errorf("only %d in stock", p.Stock)
	}
	p.Stock = p.Stock - quantity
	return nil
}
