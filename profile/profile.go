package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProfileData struct {
	ID              string           `json:"_id"`
	AccountID       string           `json:"accountId"`
	Rvn             int              `json:"rvn"`
	CommandRevision int              `json:"commandRevision"`
	Created         string           `json:"created"`
	Updated         string           `json:"updated"`
	Items           map[string]*Item `json:"items"`
	Stats           Stats            `json:"stats"`
}

type Stats struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type Item struct {
	Attributes map[string]interface{} `json:"attributes"`
}

func ReadProfile(accountId, profileId string) (*ProfileData, error) {
	path := filepath.Join("config", accountId, "profiles", profileId+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p ProfileData
	err = json.Unmarshal(data, &p)
	return &p, err
}

func SaveProfile(accountId, profileId string, data *ProfileData) error {
	path := filepath.Join("config", accountId, "profiles", profileId+".json")
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func ReadProfileTemplate(profileId string) (*ProfileData, error) {
	path := filepath.Join("config", "templates", profileId+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p ProfileData
	err = json.Unmarshal(data, &p)
	return &p, err
}

func BumpRvn(p *ProfileData) {
	p.Rvn++
	p.CommandRevision++
}

func ModifyStat(data *ProfileData, key string, value interface{}, changes *[]interface{}) {
	if data.Stats.Attributes == nil {
		data.Stats.Attributes = make(map[string]interface{})
	}
	data.Stats.Attributes[key] = value

	*changes = append(*changes, map[string]interface{}{
		"changeType": "statModified",
		"name":       key,
		"value":      value,
	})
}

func ChangeItemAttribute(data *ProfileData, itemId, key string, value interface{}, changes *[]interface{}) {
	if data.Items == nil {
		data.Items = make(map[string]*Item)
	}
	if data.Items[itemId] == nil {
		data.Items[itemId] = &Item{Attributes: make(map[string]interface{})}
	}
	data.Items[itemId].Attributes[key] = value

	*changes = append(*changes, map[string]interface{}{
		"changeType": "itemAttrChanged",
		"itemId":     itemId,
		"attribute":  key,
		"value":      value,
	})
}
