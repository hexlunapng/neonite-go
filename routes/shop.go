package routes

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func RegisterStorefrontRoutes(r *mux.Router) {

	r.HandleFunc("/fortnite/api/storefront/v2/keychain", KeychainHandler).Methods("GET")
}

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("shop.json")
	if err != nil {
		http.Error(w, `{"error":"Failed to read shop.json"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(file)
}

func KeychainHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("keychain.json") // keychain download :3
	if err != nil {
		http.Error(w, `{"error":"Failed to read keychain.json"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(file)
}
