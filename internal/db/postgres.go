package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func Connect(dbURL string) (*pgxpool.Pool, error) {
	var err error

	db, err = pgxpool.New(context.Background(), dbURL)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database not ready: %w", err)
	}

	return db, nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}
