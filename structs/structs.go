package structs

import (
	"encoding/json"
	"log"
	"net/http"
)

func NeoLog(message string) {
	log.Println("[Neonite-Go]", message)
}

type Errors struct {
	Message string `json:"message"`
}


func SendError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Errors{Message: msg})
}
