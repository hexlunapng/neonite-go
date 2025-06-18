package routes


import (
	"net/http"
	"os"
)

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("shop.json")
	if err != nil {
		http.Error(w, `{"error":"Failed to read shop.json"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(file)
}
