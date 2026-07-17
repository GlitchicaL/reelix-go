package db_test

import (
	"crypto/sha256"
	"encoding/hex"
	"reelix-go/internal/db"
	"reelix-go/internal/testdata"
	"reelix-go/internal/testutil"
	"reelix-go/internal/utils"
	"testing"
	"time"
)

// By default refresh tokens should expire after 7 days
var REFRESH_TOKEN_DAYS = 7

func TestCreateUser(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)

	if _, err := testdata.SeedUser(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	testutil.Clean()
}

func TestGetUser(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	user, _ := testdata.SeedUser()

	dbUser, err := db.GetUser(user.Username)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if user.Username != dbUser.Username {
		t.Errorf("user usernames do not match")
	}

	testutil.Clean()
}

func TestSetUserRefreshToken(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	user, _ := testdata.SeedUser()

	_, hashedRefreshToken, refreshTokenExpiresAt, err := utils.GenerateRefresh(REFRESH_TOKEN_DAYS)

	if err != nil {
		t.Errorf("error generating refresh token: %v", err)
	}

	err = db.SetUserRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiresAt)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	testutil.Clean()
}

func TestGetUserByRefreshToken(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	user, _ := testdata.SeedUser()

	refreshToken, hashedRefreshToken, refreshTokenExpiresAt, err := utils.GenerateRefresh(REFRESH_TOKEN_DAYS)

	if err != nil {
		t.Errorf("error generating refresh token: %v", err)
	}

	err = db.SetUserRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiresAt)

	if err != nil {
		t.Errorf("error setting user refresh token: %v", err)
	}

	hash := sha256.Sum256([]byte(refreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	if hashedRefreshToken != hashedToken {
		t.Errorf("hash %s does not match hash %s", hashedRefreshToken, hashedToken)
	}

	dbUser, err := db.GetUserByRefreshToken(hashedToken)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if user.Username != dbUser.Username {
		t.Errorf("user %s does not match user %s", user.Username, dbUser.Username)
	}

	if time.Now().After(dbUser.RefreshTokenExpiryDate) {
		t.Error("user refresh token expired")
	}

	testutil.Clean()
}

func TestRemoveUserRefreshToken(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	user, _ := testdata.SeedUser()

	_, hashedRefreshToken, refreshTokenExpiresAt, err := utils.GenerateRefresh(REFRESH_TOKEN_DAYS)

	if err != nil {
		t.Errorf("error generating refresh token: %v", err)
	}

	err = db.SetUserRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiresAt)

	if err != nil {
		t.Errorf("error setting user refresh token: %v", err)
	}

	err = db.RemoveUserRefreshToken(user.ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	testutil.Clean()
}

func TestGetUserCount(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	testdata.SeedUser()

	userCount, err := db.GetUserCount()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if userCount != 1 {
		t.Errorf("user count %d is not 1", userCount)
	}

	testutil.Clean()
}
