package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"reelix-go/internal/db"

	"github.com/gorilla/mux"
)

func vaultsHandler(w http.ResponseWriter, r *http.Request) {
	vaults, err := db.GetVaults()

	if err != nil {
		log.Fatalf("error fetching vaults")
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(vaults); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func vaultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vaultId, err := strconv.Atoi(vars["vaultId"])

	if err != nil {
		log.Fatalf("invalid vault id")
	}

	vault, err := db.GetVault(vaultId)

	if err != nil {
		log.Fatalf("error fetching vault")
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(vault); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func collectionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vaultId, err := strconv.Atoi(vars["vaultId"])

	if err != nil {
		log.Fatalf("invalid vault id")
	}

	collections, err := db.GetCollections(vaultId)

	if err != nil {
		log.Fatalf("error fetching collections from vault %v", vaultId)
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(collections); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func videosHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	collectionId, err := strconv.Atoi(vars["collectionId"])

	if err != nil {
		log.Fatalf("invalid collection id")
	}

	videos, err := db.GetVideos(collectionId)

	if err != nil {
		log.Fatalf("error fetching videos from collection %v", collectionId)
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(videos); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	videoId, err := strconv.Atoi(vars["videoId"])

	if err != nil {
		log.Fatalf("invalid video id")
	}

	video, err := db.GetVideo(videoId)

	if err != nil {
		log.Fatalf("error fetching video %v", videoId)
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(video); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func actorsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vaultId, err := strconv.Atoi(vars["vaultId"])

	if err != nil {
		log.Fatalf("invalid vault id")
	}

	actors, totalCount, err := db.GetActors()

	type ActorsMetadata struct {
		Actors     []db.Actor `json:"actors"`
		TotalCount int        `json:"totalCount"`
	}

	data := ActorsMetadata{
		Actors:     actors,
		TotalCount: totalCount,
	}

	if err != nil {
		log.Fatalf("error fetching actors from vault %v", vaultId)
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vaultId, err := strconv.Atoi(vars["vaultId"])

	if err != nil {
		log.Fatalf("invalid vault id")
	}

	galleries, err := db.GetGalleries(vaultId)

	if err != nil {
		log.Fatalf("error fetching galleries from vault %v", vaultId)
	}

	if err := json.NewEncoder(w).Encode(galleries); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	galleryId, err := strconv.Atoi(vars["galleryId"])

	if err != nil {
		log.Fatalf("invalid gallery id")
	}

	gallery, err := db.GetGallery(galleryId)

	if err != nil {
		log.Fatalf("error fetching gallery from vault %v", galleryId)
	}

	if err := json.NewEncoder(w).Encode(gallery); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}
