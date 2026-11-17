package main

import (
	"fmt"
	"net/http"

	"macbook-marketplace/internal/api"
	"macbook-marketplace/internal/catalog"
)

const storeName = "MacBook Marketplace"

func main() {
	fmt.Println(welcomeMessage(storeName))

	store := catalog.NewSampleStore()
	router := api.NewRouter(store)

	fmt.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Server stopped:", err)
	}
}
