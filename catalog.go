package main

// categoryExists reports whether a category with this ID exists.
func categoryExists(id int) bool {
	for _, category := range categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

// productsInCategory returns the products that belong to the category.
func productsInCategory(categoryID int) []Product {
	result := []Product{} // empty, not nil, so the JSON is [] and not null
	for _, product := range products {
		if product.CategoryID == categoryID {
			result = append(result, product)
		}
	}
	return result
}
