package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
"github.com/hex/neonite-go/routes"
	"github.com/hex/neonite-go/structs"
)

var version = "1.0"

func main() {
	port := "3551"
	structs.NeoLog("Starting server...")

	r := mux.NewRouter()
	r.Use(jsonMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Welcome to neonite`))
	})


	r.Use(jsonMiddleware)

	files, err := ioutil.ReadDir("./routes")
	if err == nil {
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".go" {
				continue
			}
		}
	}

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	structs.SendError(w, http.StatusNotFound, "not_found")
})

router := mux.NewRouter()
router.HandleFunc("/auth/v1/oauth/token", oauthTokenHandler).Methods("POST")

http.ListenAndServe(":8080", router)


	log.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, r)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
