package main

import (
	"fmt"
	"log"
	"net/http"

	"macbook-marketplace/internal/api"
	"macbook-marketplace/internal/catalog"
)

const (
	storeName = "MacBook Marketplace"
	dataFile  = "data/catalog.yml"
	address   = ":8080"
)

func main() {
	fmt.Println(welcomeMessage(storeName))

	store, err := catalog.LoadFile(dataFile)
	if err != nil {
		log.Fatal(err) // prints the error and stops: no point serving without data
	}

	router := api.NewRouter(store)

	fmt.Println("Listening on http://localhost" + address)
	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal("server stopped: ", err)
	}
}
