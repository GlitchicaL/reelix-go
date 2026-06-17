package testutil

import (
	"reelix-go/internal/db"
)

func SetupConnection(url string) error {
	_, err := db.Connect(url)

	if err != nil {
		return err
	}

	return nil
}

func Clean() {
	/*
		TODO: Could consider using transactions and rollback after each test.
	*/
	db.Exec(`
		TRUNCATE TABLE 
			vaults,
			collections,
			videos,
			tags,
			video_tags,
			actors,
			video_actors,
			galleries,
			users
		RESTART IDENTITY CASCADE
	`)

	defer db.Close()
}
