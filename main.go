package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"macbook-marketplace/internal/api"
	"macbook-marketplace/internal/catalog"
	"macbook-marketplace/internal/order"
)

const (
	storeName   = "iCycle"
	dataFile    = "data/catalog.yml"
	webDir      = "web"
	defaultPort = "8080"
)

func main() {
	fmt.Println(welcomeMessage(storeName))

	// Secrets live outside the code. godotenv copies KEY=VALUE lines from
	// the .env file into environment variables, then os.Getenv reads them.
	// If there is no .env file, real environment variables still work.
	_ = godotenv.Load()
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD is not set: copy .env.example to .env and set a password")
	}
	port := os.Getenv("PORT") // optional: most hosting services set PORT themselves
	if port == "" {
		port = defaultPort
	}
	address := ":" + port

	store, err := catalog.LoadFile(dataFile)
	if err != nil {
		log.Fatal(err) // prints the error and stops: no point serving without data
	}

	orders := order.NewService(store)
	router := api.NewRouter(store, orders, adminPassword)

	// The website: any URL that isn't an API route is a file from web/,
	// e.g. "/" → web/index.html, "/app.js" → web/app.js.
	files := http.FileServer(http.Dir(webDir))
	router.Handle("GET /", files)

	// The admin page needs the admin login too.
	router.Handle("GET /admin.html", api.AdminOnly(adminPassword, files))
	router.Handle("GET /admin.js", api.AdminOnly(adminPassword, files))

	fmt.Println("Listening on http://localhost" + address)
	fmt.Println("Admin page: http://localhost" + address + "/admin.html (user: admin)")
	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal("server stopped: ", err)
	}
}
