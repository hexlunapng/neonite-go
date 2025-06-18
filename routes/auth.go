package routes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"neonite-go/structs"

	"github.com/gorilla/mux"
)

func RegisterAccountRoutes(r *mux.Router) {
	r.HandleFunc("/account/api/oauth/token", oauthTokenHandler).Methods("POST")
	r.HandleFunc("/account/api/oauth/verify", oauthVerifyHandler).Methods("GET")
	r.HandleFunc("/account/api/oauth/sessions/kill", killSessionHandler).Methods("DELETE")
	r.HandleFunc("/account/api/oauth/sessions/kill/{token}", killSessionHandler).Methods("DELETE")

	r.HandleFunc("/account/api/public/account/{accountId}", accountByIDHandler).Methods("GET")
	r.HandleFunc("/account/api/public/account/displayName/{displayName}", accountByDisplayNameHandler).Methods("GET")
	r.HandleFunc("/account/api/public/account/", accountBatchHandler).Methods("GET")

	r.HandleFunc("/account/api/public/account/{accountId}/deviceAuth", deviceAuthListHandler).Methods("GET")
	r.HandleFunc("/account/api/public/account/{accountId}/deviceAuth", deviceAuthCreateHandler).Methods("POST")
	r.HandleFunc("/account/api/public/account/{accountId}/deviceAuth/{deviceId}", deviceAuthDeleteHandler).Methods("DELETE")
}

func oauthTokenHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		GrantType    string
		Username     string
		Code         string
		AccountID    string
		ExchangeCode string
	}

	var req Req
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			structs.SendDetailedError(w, structs.Errors["invalid_request"].With("invalid JSON"), http.StatusBadRequest)
			return
		}
	} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			structs.SendDetailedError(w, structs.Errors["invalid_request"].With("invalid form data"), http.StatusBadRequest)
			return
		}
		req.GrantType = r.FormValue("grant_type")
		req.Username = r.FormValue("username")
		req.Code = r.FormValue("code")
		req.AccountID = r.FormValue("account_id")
		req.ExchangeCode = r.FormValue("exchange_code")
	} else {
		structs.SendDetailedError(w, structs.Errors["invalid_request"].With("unsupported content type"), http.StatusUnsupportedMediaType)
		return
	}

	var displayName, accountId string

	switch req.GrantType {
	case "client_credentials", "refresh_token":
		accountId = "client_user"
		displayName = "client_user"
	case "password":
		if req.Username == "" {
			structs.SendDetailedError(w, structs.Errors["invalid_request"].With("username"), http.StatusBadRequest)
			return
		}
		displayName = strings.Split(req.Username, "@")[0]
		accountId = strings.ReplaceAll(displayName, " ", "_")
	case "authorization_code":
		if req.Code == "" {
			structs.SendDetailedError(w, structs.Errors["invalid_request"].With("code"), http.StatusBadRequest)
			return
		}
		accountId = req.Code
		displayName = req.Code
	case "device_auth":
		if req.AccountID == "" {
			structs.SendDetailedError(w, structs.Errors["invalid_request"].With("account_id"), http.StatusBadRequest)
			return
		}
		accountId = req.AccountID
		displayName = req.AccountID
	case "exchange_code":
		if req.ExchangeCode == "" {
			structs.SendDetailedError(w, structs.Errors["invalid_request"].With("exchange_code"), http.StatusBadRequest)
			return
		}
		accountId = req.ExchangeCode
		displayName = req.ExchangeCode
	default:
		structs.SendDetailedError(w, structs.Errors["unsupported_grant_type"].With(req.GrantType), http.StatusBadRequest)
		return
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		structs.SendDetailedError(w, structs.Errors["server_error"].With("failed to generate token"), http.StatusInternalServerError)
		return
	}
	accessToken := hex.EncodeToString(randomBytes)

	response := map[string]interface{}{
		"access_token":       accessToken,
		"expires_in":         28800,
		"expires_at":         "9999-12-31T23:59:59.999Z",
		"token_type":         "bearer",
		"account_id":         accountId,
		"client_id":          "ec684b8c687f479fadea3cb2ad83f5c6",
		"internal_client":    true,
		"client_service":     "fortnite",
		"refresh_token":      "STATIC_REFRESH_TOKEN",
		"refresh_expires":    115200,
		"refresh_expires_at": "9999-12-31T23:59:59.999Z",
		"displayName":        displayName,
		"app":                "fortnite",
		"in_app_id":          accountId,
		"device_id":          "static-device-id",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func oauthVerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "bearer ")

	response := map[string]interface{}{
		"access_token":       token,
		"expires_in":         28800,
		"expires_at":         "9999-12-31T23:59:59.999Z",
		"token_type":         "bearer",
		"refresh_token":      "STATIC_REFRESH_TOKEN",
		"refresh_expires":    115200,
		"refresh_expires_at": "9999-12-31T23:59:59.999Z",
		"account_id":         "ninja",
		"client_id":          "3446cd72694c4a4485d81b77adbb2141",
		"internal_client":    true,
		"client_service":     "fortnite",
		"displayName":        "ninja",
		"app":                "fortnite",
		"in_app_id":          "ninja",
		"device_id":          "static-device-id",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func killSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func accountByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["accountId"]
	response := map[string]interface{}{
		"id":            id,
		"displayName":   id,
		"externalAuths": map[string]interface{}{},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func accountByDisplayNameHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["displayName"]
	response := map[string]interface{}{
		"id":            name,
		"displayName":   name,
		"externalAuths": map[string]interface{}{},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func accountBatchHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["accountId"]
	if len(ids) == 0 {
		http.Error(w, "Missing accountId", http.StatusBadRequest)
		return
	}

	var response []map[string]interface{}
	for _, id := range ids {
		displayName := id
		if strings.HasPrefix(id, "NeoniteBot") {
			displayName = "NeoniteBot"
		}
		response = append(response, map[string]interface{}{
			"id":            id,
			"displayName":   displayName,
			"externalAuths": map[string]interface{}{},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deviceAuthListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

func deviceAuthCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := map[string]string{
		"accountId": vars["accountId"],
		"deviceId":  "null",
		"secret":    "null",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deviceAuthDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
