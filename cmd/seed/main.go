package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"reelix-go/internal/api"
	"reelix-go/internal/db"
	"reelix-go/internal/testdata"
	"reelix-go/internal/testutil"
)

func main() {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error

	for i := 1; i <= 30; i++ {
		err = testutil.SetupConnection(dbURL)

		if err == nil {
			break
		}

		log.Printf("failed to connect to database (attempt %d/30), retrying...", i)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect to database after 30 retries:", err)
	}

	defer db.Close()

	/*
		At this point the mock database is up and we can start seeding
		mock data for Postman/Playwright testing
	*/

	testdata.SeedVaults()

	router := api.NewRouter()

	fmt.Println("Reelix Test API server started on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
