package db

import (
	"context"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
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
		return nil, fmt.Errorf("error creating user: %v", err)
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
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	return &u, nil
}

func GetUserCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int

	err := db.QueryRow(
		context.Background(),
		query,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error fetching user count: %v", err)
	}

	return count, nil
}
