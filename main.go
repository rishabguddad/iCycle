package main

import "fmt"

const storeName = "MacBook Marketplace"

func main() {

	message := welcomeMessage(storeName)

	fmt.Println(message)

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
