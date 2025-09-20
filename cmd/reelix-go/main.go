package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"reelix-go/internal/api"
	"reelix-go/internal/db"
	"reelix-go/internal/scanner"
)

func main() {
	// Connect to DB
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		"database",
		"5432",
		os.Getenv("DB_NAME"),
	)

	_, err := db.Connect(dbURL)

	if err != nil {
		log.Fatal("failed to connect to database")
	}

	defer db.Close()

	scanner.Scan()

	router := api.NewRouter()

	fmt.Println("Reelix video server started on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
