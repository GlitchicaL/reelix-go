package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"reelix-go/internal/db"

	"github.com/gorilla/mux"
)

func vaultHandler(w http.ResponseWriter, r *http.Request) {
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

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	collections, err := db.GetCollections(1)

	if err != nil {
		log.Fatalf("error fetching collections from vault %v", 1)
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(collections); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func videosHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionId, err := strconv.Atoi(vars["id"])

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
	collectionId, _ := strconv.Atoi(vars["collectionId"])
	videoSlug := vars["videoSlug"]

	video, err := db.GetVideo(collectionId, videoSlug)

	if err != nil {
		log.Fatalf("error fetching videos from collection %v", collectionId)
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(video); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}
