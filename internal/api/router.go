package api

import (
	"github.com/gorilla/mux"
)

// NewRouter sets up the API routes and returns a configured router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/vaults", vaultsHandler).Methods("GET")                                  // For vaults
	r.HandleFunc("/api/vault/{vaultId}", vaultHandler).Methods("GET")                          // For vault
	r.HandleFunc("/api/collections/{vaultId}", collectionsHandler).Methods("GET")              // For collections
	r.HandleFunc("/api/galleries/{vaultId}", galleriesHandler).Methods("GET")                  // For galleries
	r.HandleFunc("/api/gallery/{galleryId}", galleryHandler).Methods("GET")                    // For gallery
	r.HandleFunc("/api/videos/{vaultId}/{collectionId}", videosHandler).Methods("GET")         // For videos
	r.HandleFunc("/api/video/{vaultId}/{collectionId}/{videoId}", videoHandler).Methods("GET") // For video
	r.HandleFunc("/api/actors/{vaultId}", actorsHandler).Methods("GET")                        // For actors

	return r
}
