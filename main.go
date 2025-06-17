
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"neonite/structs"
)

var version = "1.0"

func main() {
	port := "3551"
	structs.NeoLog("Starting server...")

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(fs)

	r.Use(jsonMiddleware)

	files, err := ioutil.ReadDir("./managers")
	if err == nil {
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".go" {
				continue
			}
		}
	}

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		structs.SendError(w, structs.Errors["not_found"], http.StatusNotFound)
	})

	log.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, r)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
