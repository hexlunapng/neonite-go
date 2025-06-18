package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegistertryPlayOnPlatformRoute(r *mux.Router) {
	r.HandleFunc("/fortnite/api/game/v2/tryPlayOnPlatform/account/{rest:.*}", TryPlayOnPlatformHandler).Methods("POST")
}

func TryPlayOnPlatformHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("true"))
}
