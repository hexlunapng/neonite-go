package main

import (
	"log"
	"net/http"

	"neonite-go/routes"
	"neonite-go/structs"

	"github.com/gorilla/mux"
)

var version = "1.0"

func main() {
	port := "3551"
	structs.NeoLog("Starting server...")

	r := mux.NewRouter()
	r.Use(jsonMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Welcome to neonite"}`))
	}).Methods("GET")

	routes.RegisterAccountRoutes(r)
	routes.RegistertryPlayOnPlatformRoute(r)
	routes.RegisterStorefrontRoutes(r)
	routes.RegisterLightswitchRoutes(r)
	routes.RegisterPermission(r)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		structs.SendError(w, http.StatusNotFound, "not_found")
	})

	addr := ":" + port
	log.Printf("Listening on port %s…\n", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
