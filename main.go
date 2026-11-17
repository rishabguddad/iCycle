package main

import (
	"fmt"
	"net/http"
)

const storeName = "MacBook Marketplace"

func main() {
	fmt.Println(welcomeMessage(storeName))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/categories", handleListCategories)

	fmt.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server stopped:", err)
	}
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
