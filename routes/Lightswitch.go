package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterLightswitchRoutes(r *mux.Router) {
	r.HandleFunc("/lightswitch/api/service/Fortnite/status", FortniteStatusHandler).Methods("GET")
	r.HandleFunc("/lightswitch/api/service/bulk/status", BulkStatusHandler).Methods("GET")
}

func FortniteStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := []map[string]interface{}{
		{
			"serviceInstanceId":  "fortnite",
			"status":             "UP",
			"message":            "Neonite-Go is UP",
			"maintenanceUri":     nil,
			"overrideCatalogIds": []string{"a7f138b2e51945ffbfdacc1af0541053"},
			"allowedActions":     []string{"PLAY", "DOWNLOAD"},
			"banned":             false,
			"launcherInfoDTO": map[string]interface{}{
				"appName":       "Fortnite",
				"catalogItemId": "4fe75bbc5a674f4f9b356b5c90567da5",
				"namespace":     "fn",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func BulkStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := []map[string]interface{}{
		{
			"serviceInstanceId":  "fortnite",
			"status":             "UP",
			"message":            "Neonite-Go is UP.",
			"maintenanceUri":     nil,
			"overrideCatalogIds": []string{"a7f138b2e51945ffbfdacc1af0541053"},
			"allowedActions":     []string{"PLAY", "DOWNLOAD"},
			"banned":             false,
			"launcherInfoDTO": map[string]interface{}{
				"appName":       "Fortnite",
				"catalogItemId": "4fe75bbc5a674f4f9b356b5c90567da5",
				"namespace":     "fn",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
