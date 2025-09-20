package api

import (
	"github.com/gorilla/mux"
)

// NewRouter sets up the API routes and returns a configured router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/vaults", vaultHandler).Methods("GET")           // For vaults
	r.HandleFunc("/api/collections", collectionHandler).Methods("GET") // For collections
	r.HandleFunc("/api/videos/{id}", videosHandler).Methods("GET")     // For videos

	r.HandleFunc("/api/video/{collectionId}/{videoSlug}", videoHandler).Methods("GET") // For videos

	return r
}
