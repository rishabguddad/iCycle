package main

import "fmt"

const storeName = "MacBook Marketplace"

func main() {
	fmt.Println(welcomeMessage(storeName))
	fmt.Println()

	printProducts()
	fmt.Println()

	trySell(1, 1)  // works: stock 1 → 0
	trySell(1, 1)  // fails: no stock left
	trySell(7, 5)  // works: stock 50 → 45
	trySell(7, 0)  // fails: invalid quantity
	trySell(99, 1) // fails: no product with ID 99
}

// printProducts prints one line per product, with its category name.
func printProducts() {
	categoryNames := make(map[int]string)
	for _, category := range categories {
		categoryNames[category.ID] = category.Name
	}

	for _, product := range products {
		categoryName, ok := categoryNames[product.CategoryID]
		if !ok {
			categoryName = "Unknown"
		}

		fmt.Printf("#%d %s | %s | %s | price: %d | stock: %d\n",
			product.ID, product.Name, categoryName,
			product.Condition, product.Price, product.Stock)
	}
}

// trySell attempts a sale and prints what happened.
func trySell(id int, quantity int) {
	product := findProduct(id)
	if product == nil {
		fmt.Printf("Product #%d not found\n", id)
		return
	}

	err := product.Sell(quantity)
	if err != nil {
		fmt.Printf("Cannot sell %d x %s: %v\n", quantity, product.Name, err)
		return
	}

	fmt.Printf("Sold %d x %s. Stock left: %d\n", quantity, product.Name, product.Stock)
}
