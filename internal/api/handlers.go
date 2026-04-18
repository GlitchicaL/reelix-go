package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	refreshTokenExp := 60 * 60 * 24 * 7 // 1 week
	accessTokenExp := 15 * time.Minute  // 15 minutes

	refreshToken, hashedRefreshToken, refreshTokenExpiresAt, err := generateRefresh(refreshTokenExp)

	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, err := generateJWT(user.Username, user.IsAdmin, accessTokenExp)

	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	if err := db.SetUserRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiresAt); err != nil {
		http.Error(w, "Error updating refresh token", http.StatusInternalServerError)
		return
	}

	log.Printf("refresh token %v", refreshToken)
	log.Printf("access token %v", accessToken)

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
		"token":    accessToken,
		"user":     user.Username,
		"is_admin": user.IsAdmin,
	})
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("reelix_refresh_token")

	if err != nil {
		http.Error(w, "Refresh token error", http.StatusInternalServerError)
		return
	}

	// This will require hashing the refresh token
	hash := sha256.Sum256([]byte(refreshCookie.Value))
	hashedToken := hex.EncodeToString(hash[:])

	user, err := db.GetUserByRefreshToken(hashedToken)

	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if time.Now().After(user.RefreshTokenExpiryDate) {
		// logout
	}

	accessTokenExp := 15 * time.Minute
	accessToken, err := generateJWT(user.Username, user.IsAdmin, accessTokenExp)

	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

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
		"token":    accessToken,
		"user":     user.Username,
		"is_admin": user.IsAdmin,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

}

func generateRefresh(tokenExp int) (token string, hashedToken string, expiresAt time.Time, err error) {
	b := make([]byte, 32) // 256-bit token
	if _, err = rand.Read(b); err != nil {
		return "", "", time.Time{}, err
	}

	token = base64.RawURLEncoding.EncodeToString(b)

	hash := sha256.Sum256([]byte(token))
	hashedToken = hex.EncodeToString(hash[:])

	expiresAt = time.Now().Add(time.Duration(tokenExp))

	return token, hashedToken, expiresAt, nil
}

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func generateJWT(username string, isAdmin bool, tokenExp time.Duration) (string, error) {
	claims := Claims{
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenExp))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")

	jwt, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return jwt, nil
}

func validateJWT(tokenStr string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
