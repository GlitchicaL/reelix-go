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

	vaultId, err := strconv.Atoi(vars["vaultId"])

	if err != nil {
		log.Fatalf("invalid vault id")
	}

	collectionId, err := strconv.Atoi(vars["collectionId"])

	if err != nil {
		log.Fatalf("invalid collection id")
	}

	videos, err := db.GetVideos(vaultId, collectionId)

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

	vaultId, _ := strconv.Atoi(vars["vaultId"])
	collectionId, _ := strconv.Atoi(vars["collectionId"])
	videoId, _ := strconv.Atoi(vars["videoId"])

	video, err := db.GetVideo(vaultId, collectionId, videoId)

	if err != nil {
		log.Fatalf("error fetching videos from collection %v", collectionId)
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(video); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func actorsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vaultId, _ := strconv.Atoi(vars["vaultId"])

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

	vaultId, _ := strconv.Atoi(vars["vaultId"])

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

	galleryId, _ := strconv.Atoi(vars["galleryId"])

	gallery, err := db.GetGallery(galleryId)

	if err != nil {
		log.Fatalf("error fetching gallery from vault %v", galleryId)
	}

	if err := json.NewEncoder(w).Encode(gallery); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}
