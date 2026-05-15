package db

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID                     int    `json:"id"`
	Username               string `json:"username"`
	Password               string `json:"password"`
	RefreshTokenHash       string
	RefreshTokenExpiryDate time.Time
	IsAdmin                bool `json:"isAdmin"`
}

func CreateUser(username string, password string, isAdmin bool) (*User, error) {
	query := `
		INSERT INTO users (username, password_hash, is_admin) VALUES ($1, $2, $3)
		RETURNING id
	`

	var u User

	err := db.QueryRow(
		context.Background(),
		query,
		username,
		password,
		isAdmin,
	).Scan(&u.ID)

	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return &u, nil
}

func GetUser(username string) (*User, error) {
	query := `SELECT id, username, password_hash, is_admin FROM users WHERE username = $1`

	var u User

	err := db.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.IsAdmin)

	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}

	return &u, nil
}

func GetUserByRefreshToken(hashedToken string) (*User, error) {
	query := `SELECT id, username, password_hash, refresh_token_hash, refresh_token_expiry_date, is_admin FROM users WHERE refresh_token_hash = $1`

	var u User

	err := db.QueryRow(
		context.Background(),
		query,
		hashedToken,
	).Scan(&u.ID, &u.Username, &u.Password, &u.RefreshTokenHash, &u.RefreshTokenExpiryDate, &u.IsAdmin)

	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}

	return &u, nil
}

func SetUserRefreshToken(id int, hashedToken string, tokenExpiryDate time.Time) error {
	query := `
		UPDATE users
		SET refresh_token_hash = $1,
			refresh_token_expiry_date = $2
		WHERE id = $3
	`

	cmdTag, err := db.Exec(
		context.Background(),
		query,
		hashedToken,
		tokenExpiryDate,
		id,
	)

	if err != nil {
		return fmt.Errorf("updating refresh token: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	return nil
}

func GetUserCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int

	err := db.QueryRow(
		context.Background(),
		query,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("fetching user count: %w", err)
	}

	return count, nil
}
