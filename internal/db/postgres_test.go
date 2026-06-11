package db

import (
	"context"
	"testing"
	"time"
)

func SetupConnection(t *testing.T) error {
	dbURL := "postgres://testuser:testpass@localhost:5433/reelix-test-db"

	db, err := Connect(dbURL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
		return err
	}

	/*
		In case the database isn't fully ready
		before unit tests are ran, let's ping
		the DB and make sure it's ready
	*/

	for i := 1; i < 5; i++ {
		err := db.Ping(context.Background())

		if err == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}

	t.Cleanup(func() {
		/*
			TODO: Could consider using transactions and rollback after each test.
		*/
		_, err := db.Exec(
			context.Background(), `
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

		if err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		defer Close()
	})

	return nil
}

func TestConnection(t *testing.T) {
	if err := SetupConnection(t); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
