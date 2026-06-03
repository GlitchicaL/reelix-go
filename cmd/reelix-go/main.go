package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	var err error

	for i := 1; i <= 30; i++ {
		_, err = db.Connect(dbURL)

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

	root := "/reelix"

	go func() {
		scan := func() {
			world, err := scanner.Scan(root)

			if err != nil {
				log.Println("error during scanning:", err)
				return
			}

			/*
				TODO: Compare the world state with a
				previous world state to determine if
				a sync is needed. Maybe also consider
				a way to only sync new changes.
			*/

			if err = scanner.Sync(world); err != nil {
				log.Println("error during syncing:", err)
			}
		}

		/*
			We'll need to run an initial scan followed
			by a scan every 15 minutes
		*/

		scan()

		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			scan()
		}
	}()

	router := api.NewRouter()

	fmt.Println("Reelix API server started on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
