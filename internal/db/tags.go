package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `xml:"tag"`
}

func CreateTags(tags []Tag) ([]Tag, error) {
	// We use make() here because at this point we know the size of
	// the slices and we won't need to reallocate memory if we were
	// to just loop and append.
	names := make([]string, len(tags))

	for i, t := range tags {
		names[i] = t.Name
	}

	query := `
		INSERT INTO tags (name) 
		SELECT *
		FROM UNNEST(
			$1::text[]
		)
		ON CONFLICT (name) DO UPDATE
		SET 
			name = EXCLUDED.name
		RETURNING id, name
	`

	rows, err := db.Query(
		context.Background(),
		query,
		names,
	)

	if err != nil {
		return nil, fmt.Errorf("creating tags: %w", err)
	}

	defer rows.Close()

	// Since we know the original size of the slice prior
	// to inserting, we know the max capacity of rows returned.
	// We don't specify length as a row conflict will result in
	// an update and not an insert.

	dbTags := make([]Tag, 0, len(tags))

	for rows.Next() {
		var t Tag

		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("creating tags: %w", err)
		}

		dbTags = append(dbTags, t)
	}

	log.Printf("tags inserted: %v", len(dbTags))

	return dbTags, nil
}

func CreateTag(tag string, tx pgx.Tx) (*int, error) {
	query := `
		INSERT INTO tags (name) VALUES ($1)
		ON CONFLICT (name) DO NOTHING
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

	log.Printf("tag inserted: %v", tag)

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
		return fmt.Errorf("linking tag %d to video: %w", tagId, err)
	}

	return nil
}

func LinkVideoTags(videos []Video, tags []Tag) error {
	/*
		Due to a video having multiple tags, we may not initially
		know the expected size of the final array.

		NOTE: Can we do len(videos) * len(tags)?
	*/

	videoIds := make([]int, 0, len(videos))
	tagsIds := make([]int, 0, len(videos))

	/*
		Since videos may have more than 1 tags, we'll have to flatten
		the videoIds array. For example, if Video [1] has Tag [5, 8]:

		videoIds = [1, 1]
		tagsIds = [5, 8]

		Since tag names are treated as unique, we can map the name to
		it's database ID
	*/

	tagsMap := map[string]int{}

	for _, t := range tags {
		tagsMap[t.Name] = t.ID
	}

	for _, v := range videos {
		if v.ID == 0 {
			continue
		}

		for _, t := range v.Tags {
			id, ok := tagsMap[t]

			if !ok {
				continue
			}

			videoIds = append(videoIds, v.ID)
			tagsIds = append(tagsIds, id)
		}
	}

	query := `
		INSERT INTO video_tags (video_id, tag_id)
		SELECT *
		FROM UNNEST(
			$1::int[],
			$2::int[]
		)
		ON CONFLICT (video_id, tag_id) DO NOTHING
	`

	rows, err := db.Query(
		context.Background(),
		query,
		videoIds,
		tagsIds,
	)

	if err != nil {
		return fmt.Errorf("linking video tags: %w", err)
	}

	defer rows.Close()

	return nil
}
