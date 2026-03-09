package main

import (
	"keep/internal/api"
	"net/http"
)

func main() {
	webMux := http.NewServeMux()
	webMux.HandleFunc("GET /keeps", api.ListKeeps)
	webMux.HandleFunc("POST /keeps/create", api.CreateKeep)
	webMux.HandleFunc("PUT /keeps/{uuid}", api.UpdateKeep)
	webMux.HandleFunc("DELETE /keeps/{uuid}", api.DeleteKeep)
	// Serve Static Frontend
	webMux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	// Start web server on 8080
	http.ListenAndServe("[::]:8080", api.CORSMiddleware(webMux))
}
