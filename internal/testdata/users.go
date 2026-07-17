package testdata

import (
	"log"
	"reelix-go/internal/db"

	"golang.org/x/crypto/bcrypt"
)

var User = db.User{
	Username: "testuser",
	Password: "testpass",
	IsAdmin:  true,
}

func SeedUser() (*db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), 14)

	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, err
	}

	dbUser, err := db.CreateUser(User.Username, string(hashedPassword), User.IsAdmin)

	if err != nil {
		return nil, err
	}

	return dbUser, nil
}
