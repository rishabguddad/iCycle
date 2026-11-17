// Package catalog holds the store's categories and products,
// and the rules for finding and selling them. It knows nothing about HTTP.
package catalog

// Categories returns every category.
func Categories() []Category {
	return categories
}

// Products returns every product.
func Products() []Product {
	return products
}

// FindProduct returns a pointer to the product with the given ID,
// or nil if there is no such product.
func FindProduct(id int) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

// CategoryExists reports whether a category with this ID exists.
func CategoryExists(id int) bool {
	for _, category := range categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

// ProductsInCategory returns the products that belong to the category.
func ProductsInCategory(categoryID int) []Product {
	result := []Product{} // empty, not nil, so the JSON is [] and not null
	for _, product := range products {
		if product.CategoryID == categoryID {
			result = append(result, product)
		}
	}
	return result
}
