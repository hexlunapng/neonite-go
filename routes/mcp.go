package routes

import (
	"encoding/json"
	"fmt"
	"neonite-go/structs"
	"neonite-go/routes"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type CommandRequest struct {
	SourceIndex             int      `json:"sourceIndex"`
	TargetIndex             int      `json:"targetIndex"`
	OptNewNameForTarget     string   `json:"optNewNameForTarget"`
	Name                    string   `json:"name"`
	LockerItem              string   `json:"lockerItem"`
	ItemIds                 []string `json:"itemIds"`
	OfferId                 string   `json:"offerId"`
	GiftBoxItemIds          []string `json:"giftBoxItemIds"`
	AffiliateName           string   `json:"affiliateName"`
	SlotName                string   `json:"slotName"`
	IndexWithinSlot         int      `json:"indexWithinSlot"`
	ItemToSlot              string   `json:"itemToSlot"`
	NewPlatform             string   `json:"newPlatform"`
	BReceiveGifts           bool     `json:"bReceiveGifts"`
	BFavorite               bool     `json:"bFavorite"`
	TargetItemId            string   `json:"targetItemId"`
	ItemFavStatus           []bool   `json:"itemFavStatus"`
	Archived                bool     `json:"archived"`
}

func ProfileCommandHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	accountId := chi.URLParam(r, "accountId")
	command := chi.URLParam(r, "command")
	profileId := r.URL.Query().Get("profileId")
	if profileId == "" {
		profileId = "common_core"
	}

	getOrCreateProfile := func(profileId string) (*profile.ProfileData, *structs.ProfileResponse, error) {
		data, err := profile.ReadProfile(accountId, profileId)
		if err != nil || data == nil {
			tmpl, err := profile.ReadProfileTemplate(profileId)
			if err != nil || tmpl == nil {
				return nil, nil, structs.NewAPIError("operation_forbidden").With(profileId)
			}
			tmpl.Created = time.Now().UTC().Format(time.RFC3339)
			tmpl.Updated = tmpl.Created
			tmpl.AccountID = accountId
			tmpl.ID = accountId

			dir := filepath.Join("config", accountId, "profiles")
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, nil, fmt.Errorf("failed to create profile directory: %w", err)
			}
			if err := profile.SaveProfile(accountId, profileId, tmpl); err != nil {
				return nil, nil, err
			}
			data = tmpl
		}

		resp := &structs.ProfileResponse{
			ProfileRevision:            data.Rvn,
			ProfileId:                  profileId,
			ProfileChangesBaseRevision: data.Rvn,
			ProfileCommandRevision:     data.CommandRevision,
			ResponseVersion:            1,
			ServerTime:                 time.Now().UTC().Format(time.RFC3339),
		}
		return data, resp, nil
	}

	data, response, err := getOrCreateProfile(profileId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	var reqBody CommandRequest
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	switch command {
	case "CopyCosmeticLoadout":
		if profileId != "athena" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		if reqBody.SourceIndex == 0 {
			data.Items[fmt.Sprintf("neoset%d_loadout", reqBody.TargetIndex)] = data.Items["sandbox_loadout"]
			data.Items[fmt.Sprintf("neoset%d_loadout", reqBody.TargetIndex)].Attributes["locker_name"] = reqBody.OptNewNameForTarget
			data.Stats.Attributes["loadouts"].([]string)[reqBody.TargetIndex] = fmt.Sprintf("neoset%d_loadout", reqBody.TargetIndex)
		} else {
			item := data.Items[fmt.Sprintf("neoset%d_loadout", reqBody.SourceIndex)]
			if item == nil {
				utils.WriteError(w, structs.NewAPIError("item_not_found").With(reqBody.LockerItem))
				return
			}
			data.Stats.Attributes["active_loadout_index"] = reqBody.SourceIndex
			data.Stats.Attributes["last_applied_loadout"] = fmt.Sprintf("neoset%d_loadout", reqBody.SourceIndex)
			data.Items["sandbox_loadout"].Attributes["locker_slots_data"] = item.Attributes["locker_slots_data"]
		}
	case "DeleteCosmeticLoadout":
		if profileId != "athena" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		loadouts := data.Stats.Attributes["loadouts"].([]string)
		loadouts[reqBody.TargetIndex] = ""
		data.Stats.Attributes["loadouts"] = loadouts
	case "SetMtxPlatform":
		if profileId != "common_core" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		response.ProfileChanges = append(response.ProfileChanges, structs.ProfileChange{
			ChangeType: "statModified",
			Name:       "current_mtx_platform",
			Value:      reqBody.NewPlatform,
		})
	case "SetReceiveGiftsEnabled":
		if profileId != "common_core" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		profile.ModifyStat(data, "allowed_to_receive_gifts", reqBody.BReceiveGifts, &response.ProfileChanges)
	case "SetItemFavoriteStatus":
		if profileId != "athena" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		item := data.Items[reqBody.TargetItemId]
		if item != nil && item.Attributes["favorite"] != reqBody.BFavorite {
			profile.ChangeItemAttribute(data, reqBody.TargetItemId, "favorite", reqBody.BFavorite, &response.ProfileChanges)
		}
	case "SetItemFavoriteStatusBatch":
		if profileId != "athena" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		for i, itemId := range reqBody.ItemIds {
			if i < len(reqBody.ItemFavStatus) {
				profile.ChangeItemAttribute(data, itemId, "favorite", reqBody.ItemFavStatus[i], &response.ProfileChanges)
			}
		}
	case "SetItemArchivedStatusBatch":
		if profileId != "athena" {
			utils.WriteError(w, structs.NewAPIError("invalid_profile").With(profileId))
			return
		}
		for _, itemId := range reqBody.ItemIds {
			profile.ChangeItemAttribute(data, itemId, "archived", reqBody.Archived, &response.ProfileChanges)
		}
	default:
		utils.WriteError(w, structs.NewAPIError("unsupported_command").With(command))
		return
	}

	if len(response.ProfileChanges) > 0 {
		profile.BumpRvn(data)
		response.ProfileRevision = data.Rvn
		response.ProfileCommandRevision = data.CommandRevision
		profile.SaveProfile(accountId, profileId, data)
	}

 json.NewEncoder(w).Encode(response)
}
