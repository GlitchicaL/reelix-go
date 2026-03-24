package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"reelix-go/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type StatusMetadata struct {
	Status string `json:"status"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	data := StatusMetadata{
		Status: "OK",
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

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

	vault, err := db.GetVault(vaultId)

	if err != nil {
		log.Fatalf("error fetching vault %v", vaultId)
	}

	actors, err := db.GetActors(vaultId)

	type ActorsMetadata struct {
		Actors    []db.Actor `json:"actors"`
		VaultName string     `json:"vaultName"`
	}

	data := ActorsMetadata{
		Actors:    actors,
		VaultName: vault.Name,
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	_, err := db.GetUser(input.Username)

	if err == nil {
		http.Error(w, "Username already in use", http.StatusConflict)
		return
	}

	count, err := db.GetUserCount()

	if err != nil {
		http.Error(w, "Error fetching user count", http.StatusInternalServerError)
		return
	}

	var isAdmin = false

	// Our first registered user will be admin by default
	if count == 0 {
		isAdmin = true
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.CreateUser(input.Username, string(hashedPassword), isAdmin)

	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := db.GetUser(input.Username)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(user.Username, user.IsAdmin)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    token,
		"user":     user.Username,
		"is_admin": user.IsAdmin,
	})
}

func generateJWT(username string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	jwt, _ := token.SignedString([]byte("your-secret-jwt-key"))

	return jwt, nil
}
