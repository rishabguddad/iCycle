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
	mux.HandleFunc("GET /api/categories/{id}/products", handleListCategoryProducts)
	mux.HandleFunc("GET /api/products", handleListProducts)
	mux.HandleFunc("GET /api/products/{id}", handleGetProduct)

	fmt.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server stopped:", err)
	}
}
