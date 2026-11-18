package main

import (
	"fmt"
	"log"
	"net/http"

	"macbook-marketplace/internal/api"
	"macbook-marketplace/internal/catalog"
	"macbook-marketplace/internal/order"
)

const (
	storeName = "MacBook Marketplace"
	dataFile  = "data/catalog.yml"
	webDir    = "web"
	address   = ":8080"
)

func main() {
	fmt.Println(welcomeMessage(storeName))

	store, err := catalog.LoadFile(dataFile)
	if err != nil {
		log.Fatal(err) // prints the error and stops: no point serving without data
	}

	orders := order.NewService(store)
	router := api.NewRouter(store, orders)

	// The website: any URL that isn't an API route is a file from web/,
	// e.g. "/" → web/index.html, "/app.js" → web/app.js.
	router.Handle("GET /", http.FileServer(http.Dir(webDir)))

	fmt.Println("Listening on http://localhost" + address)
	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal("server stopped: ", err)
	}
}
