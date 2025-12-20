package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func CreateTag(tag string, tx pgx.Tx) (*int, error) {
	query := `
		INSERT INTO tags (name) VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`

	var tagId int

	err := tx.QueryRow(
		context.Background(),
		query,
		tag,
	).Scan(&tagId)

	if err != nil {
		return nil, fmt.Errorf("failed to insert tag %s: %w", tag, err)
	}

	log.Printf("tag added: %v", tag)

	return &tagId, nil
}

func LinkVideoTag(videoId int, tagId int, tx pgx.Tx) error {
	query := `
			INSERT INTO video_tags (video_id, tag_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

	_, err := tx.Exec(
		context.Background(),
		query,
		videoId,
		tagId,
	)

	if err != nil {
		return fmt.Errorf("failed to link tag %v to video: %w", tagId, err)
	}

	return nil
}
