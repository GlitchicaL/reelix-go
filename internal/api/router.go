package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/status", statusHandler).Methods("GET")

	r.HandleFunc("/api/vaults", vaultsHandler).Methods("GET")
	r.HandleFunc("/api/vault/{vaultId}", vaultHandler).Methods("GET")

	r.HandleFunc("/api/collections/{vaultId}", collectionsHandler).Methods("GET")

	r.HandleFunc("/api/videos/{collectionId}", videosHandler).Methods("GET")
	r.HandleFunc("/api/video/{videoId}", videoHandler).Methods("GET")

	r.HandleFunc("/api/galleries/{vaultId}", galleriesHandler).Methods("GET")
	r.HandleFunc("/api/gallery/{galleryId}", galleryHandler).Methods("GET")

	r.HandleFunc("/api/actors/{vaultId}", actorsHandler).Methods("GET")

	return r
}
