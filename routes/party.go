package party

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type MemberConnectionMeta struct {
	Platform string `json:"urn:epic:conn:platform_s"`
	Type     string `json:"urn:epic:conn:type_s"`
}

type MemberConnection struct {
	ID          string            `json:"id"`
	ConnectedAt time.Time         `json:"connected_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	YieldLead   bool              `json:"yield_leadership"`
	Meta        MemberConnectionMeta `json:"meta"`
}

type MemberMeta struct {
	DisplayName string `json:"urn:epic:member:dn_s,omitempty"`
}

type Member struct {
	AccountID   string             `json:"account_id"`
	Meta        MemberMeta         `json:"meta"`
	Connections []MemberConnection `json:"connections"`
	Revision    int                `json:"revision"`
	UpdatedAt   time.Time          `json:"updated_at"`
	JoinedAt    time.Time          `json:"joined_at"`
	Role        string             `json:"role"`
}

type PartyConfig struct {
	Type            string `json:"type"`
	Joinability     string `json:"joinability"`
	Discoverability string `json:"discoverability"`
	SubType         string `json:"sub_type"`
	MaxSize         int    `json:"max_size"`
	InviteTTL       int    `json:"invite_ttl"`
	JoinConfirmation bool  `json:"join_confirmation"`
}

type PartyResponse struct {
	ID         string      `json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Config     PartyConfig `json:"config"`
	Members    []Member    `json:"members"`
	Applicants []string    `json:"applicants"`
	Meta       interface{} `json:"meta"`
	Invites    []string    `json:"invites"`
	Revision   int         `json:"revision"`
}

type JoinInfoConnection struct {
	ID   string                 `json:"id"`
	Meta map[string]interface{} `json:"meta"`
}

type JoinInfo struct {
	Connection JoinInfoConnection `json:"connection"`
}

type PartyRequest struct {
	Config   PartyConfig `json:"config"`
	Meta     interface{} `json:"meta"`
	JoinInfo JoinInfo    `json:"join_info"`
}

func RegisterRoutes(r *mux.Router) {
	base := "/party/api/v1"
	r.HandleFunc(base+"/{any:.*}/parties", CreateParty).Methods("POST")
	r.HandleFunc(base+"/{any:.*}/parties/{partyId}/members/{accountId}/meta", PatchMemberMeta).Methods("PATCH")
	r.HandleFunc(base+"/{any:.*}/user/{accountId}/pings/{pingerId}", PostUserPing).Methods("POST")
	r.HandleFunc(base+"/{any:.*}/parties/{partyId}/members/{accountId}/promote", EmptyHandler).Methods("POST")
	r.HandleFunc(base+"/{any:.*}/parties/{partyId}/members/{accountId}/confirm", ForbiddenHandler).Methods("POST")
	r.HandleFunc(base+"/{any:.*}/parties/{partyId}/members/{accountId}", DeleteMember).Methods("DELETE")
	r.HandleFunc(base+"/{any:.*}/user/{accountId}", GetUserParty).Methods("GET")
	r.HandleFunc(base+"/{any:.*}/parties/{partyId}", GetParty).Methods("GET")
	r.HandleFunc(base+"/{any:.*}", EmptyListHandler).Methods("ALL")
}

func CreateParty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PartyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := PartyResponse{
			ID:        "LobbyBotPartyLMFAO",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Config:    req.Config,
			Members:   []Member{},
			Applicants: []string{},
			Meta:      req.Meta,
			Invites:   []string{},
			Revision:  0,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	accountID := ""
	meta := MemberMeta{}
	connections := []MemberConnection{}

	if req.JoinInfo.Connection.ID != "" {
		accountID = strings.Split(req.JoinInfo.Connection.ID, "@")[0]
		connections = []MemberConnection{{
			ID:          req.JoinInfo.Connection.ID,
			ConnectedAt: time.Now(),
			UpdatedAt:   time.Now(),
			YieldLead:   false,
			Meta: MemberConnectionMeta{
				Platform: "WIN",
				Type:     "game",
			},
		}}
		meta.DisplayName = accountID
	}

	members := []Member{
		{
			AccountID:   accountID,
			Meta:        meta,
			Connections: connections,
			Revision:    0,
			UpdatedAt:   time.Now(),
			JoinedAt:    time.Now(),
			Role:        "CAPTAIN",
		},
	}

	resp := PartyResponse{
		ID:         "LobbyBotPartyLMFAO",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Config:     req.Config,
		Members:    members,
		Applicants: []string{},
		Meta:       req.Meta,
		Invites:    []string{},
		Revision:   0,
	}

	json.NewEncoder(w).Encode(resp)
}

func PatchMemberMeta(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func PostUserPing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]
	pingerId := vars["pingerId"]

	resp := map[string]interface{}{
		"sent_by":   pingerId,
		"sent_to":   accountId,
		"sent_at":   time.Now(),
		"expires_at": time.Now().Add(1 * time.Hour),
		"meta":      map[string]interface{}{},
	}

	json.NewEncoder(w).Encode(resp)

}

func GetUserParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	resp := map[string]interface{}{
		"current": []PartyResponse{
			{
				ID:        "LobbyBotPartyLMFAO",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Config: PartyConfig{
					Type:             "DEFAULT",
					Joinability:      "INVITE_AND_FORMER",
					Discoverability:  "INVITED_ONLY",
					SubType:          "default",
					MaxSize:          16,
					InviteTTL:        14400,
					JoinConfirmation: true,
				},
				Members: []Member{
					{
						AccountID: accountId,
						Meta: MemberMeta{
							DisplayName: accountId,
						},
						Connections: []MemberConnection{
							{
								ID:          "",
								ConnectedAt: time.Now(),
								UpdatedAt:   time.Now(),
								YieldLead:   false,
								Meta: MemberConnectionMeta{
									Platform: "WIN",
									Type:     "game",
								},
							},
						},
						Revision:  0,
						UpdatedAt: time.Now(),
						JoinedAt:  time.Now(),
						Role:      "CAPTAIN",
					},
				},
				Applicants: []string{},
				Meta:       map[string]interface{}{},
				Invites:    []string{},
				Revision:   0,
			},
		},
		"pending": []interface{}{},
		"invites": []interface{}{},
		"pings":   []interface{}{},
	}

	json.NewEncoder(w).Encode(resp)
}

func GetParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	resp := PartyResponse{
		ID:        vars["partyId"],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config: PartyConfig{
			Type:             "DEFAULT",
			Joinability:      "OPEN",
			Discoverability:  "ALL",
			SubType:          "default",
			MaxSize:          16,
			InviteTTL:        14400,
			JoinConfirmation: false,
		},
		Members:    []Member{},
		Applicants: []string{},
		Meta:       map[string]interface{}{},
		Invites:    []string{},
		Revision:   0,
	}

	json.NewEncoder(w).Encode(resp)
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}



func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func EmptyListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("[]"))
}
