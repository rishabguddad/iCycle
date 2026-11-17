package main

import (
	"encoding/json"
	"net/http"
)

// handleHealth answers GET /api/health, so anyone can check the server is running.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleListCategories answers GET /api/categories with every category.
func handleListCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, categories)
}

// writeJSON sends data to the client as JSON, with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
