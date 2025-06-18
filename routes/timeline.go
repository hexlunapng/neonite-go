package routes


import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"context"

	"github.com/go-chi/chi/v5"
)

type VersionResponse struct {
	Version string `json:"version"`
}

type CalendarResponse struct {
	Raw json.RawMessage
}

func main() {
	r := chi.NewRouter()
	r.Get("/fortnite/api/calendar/v1/timeline", TimelineHandler)
	http.ListenAndServe(":8080", r)
}

func TimelineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	version := "1.0"
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "" {
		parts := strings.Split(userAgent, "-")
		if len(parts) > 1 {
			subparts := strings.Split(parts[1], "-CL")
			if len(subparts) > 0 {
				version = subparts[0]
			}
		}
	}

	currentVersion, err := fetchCurrentVersion(ctx)
	if err != nil {
		http.Error(w, "Failed to get current version", http.StatusInternalServerError)
		return
	}

	if currentVersion == version {
		calendarData, err := fetchCalendar(ctx)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(calendarData)
			return
		}
	}

	season := "1"
	if dotIndex := strings.Index(version, "."); dotIndex != -1 {
		season = version[:dotIndex]
	} else {
		season = version
	}

	resp := map[string]interface{}{
		"channels": map[string]interface{}{
			"standalone-store":    map[string]interface{}{},
			"client-matchmaking":  map[string]interface{}{},
			"tk":                  map[string]interface{}{},
			"featured-islands":    map[string]interface{}{},
			"community-votes":     map[string]interface{}{},
			"client-events": map[string]interface{}{
				"states": []map[string]interface{}{
					{
						"validFrom": "2020-05-21T18:36:38.383Z",
						"activeEvents": []map[string]interface{}{
							{
								"eventType":   "EventFlag.LobbySeason" + season,
								"activeUntil": "9999-12-31T23:59:59.999Z",
								"activeSince": "2019-12-31T23:59:59.999Z",
							},
						},
						"state": map[string]interface{}{
							"activeStorefronts":       []interface{}{},
							"eventNamedWeights":       map[string]interface{}{},
							"activeEvents":            []interface{}{},
							"seasonNumber":            parseInt(season),
							"seasonTemplateId":        "AthenaSeason:athenaseason" + season,
							"matchXpBonusPoints":      0,
							"eventPunchCardTemplateId": "",
							"seasonBegin":             "9999-12-31T23:59:59.999Z",
							"seasonEnd":               "9999-12-31T23:59:59.999Z",
							"seasonDisplayedEnd":      "9999-12-31T23:59:59.999Z",
							"weeklyStoreEnd":          "9999-12-31T23:59:59.999Z",
							"stwEventStoreEnd":        "9999-12-31T23:59:59.999Z",
							"stwWeeklyStoreEnd":       "9999-12-31T23:59:59.999Z",
							"dailyStoreEnd":           "9999-12-31T23:59:59.999Z",
						},
					},
				},
				"cacheExpire": "9999-12-31T23:59:59.999Z",
			},
		},
		"cacheIntervalMins": 99999,
		"currentTime":       time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func fetchCurrentVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://fortnite-public-service-prod.ak.epicgames.com/fortnite/api/version", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versionResp VersionResponse
	if err := json.Unmarshal(body, &versionResp); err != nil {
		return "", err
	}

	return versionResp.Version, nil
}

func fetchCalendar(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.nitestats.com/v1/epic/modes", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func parseInt(s string) int {
	i := 1
	for _, c := range s {
		if c >= '0' && c <= '9' {
			i = i*10 + int(c-'0')
		} else {
			break
		}
	}
	if i == 1 && s != "" {
		return 0
	}
	return i
}
