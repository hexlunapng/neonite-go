package managers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"neonite/structs"
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
		GrantType    string `json:"grant_type"`
		Username     string `json:"username"`
		Code         string `json:"code"`
		AccountID    string `json:"account_id"`
		ExchangeCode string `json:"exchange_code"`
	}

	var req Req
	json.NewDecoder(r.Body).Decode(&req)

	var displayName, accountId string

	switch req.GrantType {
	case "client_credentials", "refresh_token":

	case "password":
		if req.Username == "" {
			structs.SendError(w, structs.Errors["invalid_request"].With("username"), http.StatusBadRequest)
			return
		}
		if strings.Contains(req.Username, "@") {
			displayName = strings.Split(req.Username, "@")[0]
		} else {
			displayName = req.Username
		}
		accountId = strings.ReplaceAll(displayName, " ", "_")
	case "authorization_code":
		if req.Code == "" {
			structs.SendError(w, structs.Errors["invalid_request"].With("code"), http.StatusBadRequest)
			return
		}
		displayName = req.Code
		accountId = req.Code
	case "device_auth":
		if req.AccountID == "" {
			structs.SendError(w, structs.Errors["invalid_request"].With("account_id"), http.StatusBadRequest)
			return
		}
		displayName = req.AccountID
		accountId = req.AccountID
	case "exchange_code":
		if req.ExchangeCode == "" {
			structs.SendError(w, structs.Errors["invalid_request"].With("exchange_code"), http.StatusBadRequest)
			return
		}
		displayName = req.ExchangeCode
		accountId = req.ExchangeCode
	default:
		structs.SendError(w, structs.Errors["unsupported_grant_type"].With(req.GrantType), http.StatusBadRequest)
		return
	}

	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	accessToken := hex.EncodeToString(randomBytes)

	response := map[string]interface{}{
		"access_token":          accessToken,
		"expires_in":            28800,
		"expires_at":            "9999-12-31T23:59:59.999Z",
		"token_type":            "bearer",
		"account_id":            accountId,
		"client_id":             "ec684b8c687f479fadea3cb2ad83f5c6",
		"internal_client":       true,
		"client_service":        "fortnite",
		"refresh_token":         "STATIC_REFRESH_TOKEN",
		"refresh_expires":       115200,
		"refresh_expires_at":    "9999-12-31T23:59:59.999Z",
		"displayName":           displayName,
		"app":                   "fortnite",
		"in_app_id":             accountId,
		"device_id":             "static-device-id",
	}
	json.NewEncoder(w).Encode(response)
}

func oauthVerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "bearer ")

	response := map[string]interface{}{
		"access_token":         token,
		"expires_in":           28800,
		"expires_at":           "9999-12-31T23:59:59.999Z",
		"token_type":           "bearer",
		"refresh_token":        "STATIC_REFRESH_TOKEN",
		"refresh_expires":      115200,
		"refresh_expires_at":   "9999-12-31T23:59:59.999Z",
		"account_id":           "ninja",
		"client_id":            "3446cd72694c4a4485d81b77adbb2141",
		"internal_client":      true,
		"client_service":       "fortnite",
		"displayName":          "ninja",
		"app":                  "fortnite",
		"in_app_id":            "ninja",
		"device_id":            "static-device-id",
	}
	json.NewEncoder(w).Encode(response)
}

func killSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func accountByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["accountId"]
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            id,
		"displayName":   id,
		"externalAuths": map[string]interface{}{},
	})
}

func accountByDisplayNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["displayName"]
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            name,
		"displayName":   name,
		"externalAuths": map[string]interface{}{},
	})
}

func accountBatchHandler(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["accountId"]
	if !ok || len(ids) == 0 {
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

	json.NewEncoder(w).Encode(response)
}

func deviceAuthListHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]interface{}{})
}

func deviceAuthCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(map[string]string{
		"accountId": vars["accountId"],
		"deviceId":  "null",
		"secret":    "null",
	})
}

func deviceAuthDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
