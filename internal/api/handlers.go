package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"reelix-go/internal/db"
	"reelix-go/internal/utils"

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
		log.Printf("unable to encode metadata: %v", err)
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func vaultsHandler(w http.ResponseWriter, r *http.Request) {
	vaults, err := db.GetVaults()

	if err != nil {
		log.Printf("error fetching vaults: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
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
		log.Printf("invalid vault id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	vault, err := db.GetVault(vaultId)

	if err != nil {
		log.Printf("error fetching vault %d: %v", vaultId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
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
		log.Printf("invalid vault id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	collections, err := db.GetCollections(vaultId)

	if err != nil {
		log.Printf("error fetching collections from vault %d: %v", vaultId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(collections); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	collectionId, err := strconv.Atoi(vars["collectionId"])

	if err != nil {
		log.Printf("invalid collection id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	vault, err := db.GetCollection(collectionId)

	if err != nil {
		log.Printf("error fetching collection %d: %v", collectionId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(vault); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func videosHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	collectionId, err := strconv.Atoi(vars["collectionId"])

	if err != nil {
		log.Printf("invalid collection id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	videos, err := db.GetVideos(collectionId)

	if err != nil {
		log.Printf("error fetching videos from collection %d: %v", collectionId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
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
		log.Printf("invalid video id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	video, err := db.GetVideo(videoId)

	if err != nil {
		log.Printf("error fetching video %d: %v", videoId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
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
		log.Printf("invalid vault id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	vault, err := db.GetVault(vaultId)

	if err != nil {
		log.Printf("error fetching vault %d: %v", vaultId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	actors, err := db.GetActors(vaultId)

	type ActorsMetadata struct {
		Actors []db.Actor `json:"actors"`
		Vault  db.Vault   `json:"vault"`
	}

	data := ActorsMetadata{
		Actors: actors,
		Vault:  *vault,
	}

	if err != nil {
		log.Printf("error fetching actors from vault %d: %v", vaultId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func actorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	actorId, err := strconv.Atoi(vars["actorId"])

	if err != nil {
		log.Printf("invalid actor id: %v", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	actor, err := db.GetActor(actorId)

	if err != nil {
		log.Printf("error fetching actor %d: %v", actorId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	// Respond with the metadata as JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(actor); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vaultId, err := strconv.Atoi(vars["vaultId"])

	if err != nil {
		log.Printf("invalid vault id")
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	galleries, err := db.GetGalleries(vaultId)

	if err != nil {
		log.Printf("error fetching galleries from vault %d: %v", vaultId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(galleries); err != nil {
		http.Error(w, "Unable to encode metadata", http.StatusInternalServerError)
	}
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	galleryId, err := strconv.Atoi(vars["galleryId"])

	if err != nil {
		log.Printf("invalid gallery id")
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	gallery, err := db.GetGallery(galleryId)

	if err != nil {
		log.Printf("error fetching gallery %d: %v", galleryId, err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
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
		log.Printf("invalid json: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	_, err := db.GetUser(input.Username)

	if err == nil {
		log.Printf("user already exists")
		http.Error(w, "Username already in use", http.StatusConflict)
		return
	}

	count, err := db.GetUserCount()

	if err != nil {
		log.Printf("error fetching user count: %v", err)
		http.Error(w, "Error fetching user count", http.StatusInternalServerError)
		return
	}

	var isAdmin = false

	/*
		Our first registered user will be admin by default
	*/

	if count == 0 {
		isAdmin = true
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	if err != nil {
		log.Printf("error hashing password: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.CreateUser(input.Username, string(hashedPassword), isAdmin)

	if err != nil {
		log.Printf("error creating user: %v", err)
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
		log.Printf("invalid json: %v", err)
		http.Error(w, "Invalid JSON", http.StatusUnauthorized)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(input.Username)

	if err != nil {
		log.Printf("invalid credentials: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("invalid credentials: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	refreshTokenDays := 7 // 7 days
	refreshTokenExp := 60 * 60 * 24 * refreshTokenDays
	accessTokenExp := 15 * time.Minute

	refreshToken, hashedRefreshToken, refreshTokenExpiresAt, err := utils.GenerateRefresh(refreshTokenDays)

	if err != nil {
		log.Printf("error generating refresh token: %v", err)
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, err := utils.GenerateJWT(user.Username, user.IsAdmin, accessTokenExp)

	if err != nil {
		log.Printf("error generating access token: %v", err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	if err := db.SetUserRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiresAt); err != nil {
		log.Printf("error updating refresh token: %v", err)
		http.Error(w, "Error updating refresh token", http.StatusInternalServerError)
		return
	}

	/*
		Even though the access cookie has same age as the
		refresh cookie, the JWT expiry claim is set to 15 minutes.
		This is to just avoid cases where the frontend has no access
		cookie and will allow for silent refreshing.
	*/

	refreshCookie := &http.Cookie{
		Name:     "reelix_refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   refreshTokenExp,
	}

	accessCookie := &http.Cookie{
		Name:     "reelix_auth_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   refreshTokenExp,
	}

	w.Header().Set("Content-Type", "application/json")

	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":     user.Username,
		"is_admin": user.IsAdmin,
	})
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("reelix_refresh_token")

	if err != nil {
		log.Printf("refresh token not passed: %v", err)
		http.Error(w, "Refresh token error", http.StatusUnauthorized)
		return
	}

	// This will require hashing the refresh token
	hash := sha256.Sum256([]byte(refreshCookie.Value))
	hashedToken := hex.EncodeToString(hash[:])

	user, err := db.GetUserByRefreshToken(hashedToken)

	if err != nil {
		log.Printf("invalid refresh token: %v", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if time.Now().After(user.RefreshTokenExpiryDate) {
		log.Printf("refresh token expired: %v", err)
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	accessTokenExp := 15 * time.Minute
	accessToken, err := utils.GenerateJWT(user.Username, user.IsAdmin, accessTokenExp)

	if err != nil {
		log.Printf("error generating access token: %v", err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	/*
		Even though the access cookie has same age as the
		refresh cookie, the JWT expiry claim is set to 15 minutes.
		This is to just avoid cases where the frontend has no access
		cookie and will allow for silent refreshing.
	*/

	accessCookie := &http.Cookie{
		Name:     "reelix_auth_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   refreshCookie.MaxAge,
	}

	w.Header().Set("Content-Type", "application/json")

	http.SetCookie(w, accessCookie)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":     user.Username,
		"is_admin": user.IsAdmin,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("reelix_refresh_token")

	if err != nil {
		log.Printf("refresh token not passed: %v", err)
		http.Error(w, "Refresh token error", http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(refreshCookie.Value))
	hashedToken := hex.EncodeToString(hash[:])

	user, err := db.GetUserByRefreshToken(hashedToken)

	if err != nil {
		log.Printf("invalid refresh token: %v", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	err = db.RemoveUserRefreshToken(user.ID)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "reelix_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "reelix_auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   -1,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":     user.Username,
		"is_admin": user.IsAdmin,
	})
}
