package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterPermission(r *mux.Router) {
	r.HandleFunc("/fortnite/api/game/v2/grant_access/{rest:.*}", GrantAccessHandler).Methods("POST")
	r.HandleFunc("/waitingroom/api/waitingroom", WaitingRoomHandler).Methods("GET")
	r.HandleFunc("/fortnite/api/game/v2/enabled_features", EnabledFeaturesHandler).Methods("GET")
}

func GrantAccessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`{}`))
}

func WaitingRoomHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func EnabledFeaturesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}
