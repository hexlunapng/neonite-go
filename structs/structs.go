package structs

import (
	"encoding/json"
	"log"
	"net/http"
)

func NeoLog(message string) {
	log.Println("[Neonite-Go]", message)
}

func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

type ProfileResponse struct {
	ProfileRevision            int
	ProfileId                  string
	ProfileChangesBaseRevision int
	ProfileCommandRevision     int
	ResponseVersion            int
	ServerTime                 string
	ProfileChanges             []ProfileChange
}

type ProfileChange struct {
	ChangeType string
	Name       string
	Value      interface{}
}
