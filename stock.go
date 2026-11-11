package main

// findProduct returns a pointer to the product with the given ID,
// or nil if there is no such product.
func findProduct(id int) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

// sell reduces the product's stock by quantity.
// It returns false, and changes nothing, if the quantity is invalid
// or there isn't enough stock.
func sell(product *Product, quantity int) bool {
	if quantity <= 0 || quantity > product.Stock {
		return false
	}
	product.Stock = product.Stock - quantity
	return true
}
