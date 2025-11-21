package api

import (
	"github.com/gorilla/mux"
)

// NewRouter sets up the API routes and returns a configured router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/vaults", vaultHandler).Methods("GET")                                   // For vaults
	r.HandleFunc("/api/collections/{vaultId}", collectionHandler).Methods("GET")               // For collections
	r.HandleFunc("/api/videos/{vaultId}/{collectionId}", videosHandler).Methods("GET")         // For videos
	r.HandleFunc("/api/video/{vaultId}/{collectionId}/{videoId}", videoHandler).Methods("GET") // For videos

	return r
}
