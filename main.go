package main

import (
	"fmt"
	"net/http"

	"macbook-marketplace/internal/api"
)

const storeName = "MacBook Marketplace"

func main() {
	fmt.Println(welcomeMessage(storeName))

	router := api.NewRouter()

	fmt.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Server stopped:", err)
	}
}
