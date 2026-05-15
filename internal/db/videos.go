package db

import (
	"context"
	"fmt"
)

type Video struct {
	ID             int      `json:"id"`
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	Studio         string   `json:"studio"`
	Tags           []string `json:"tags"`
	Actors         []Actor  `json:"actors"`
	CollectionID   int      `json:"collectionId"`
	CollectionName string   `json:"collectionName"`
	CollectionSlug string   `json:"collectionSlug"`
	VaultID        int      `json:"vaultId"`
	VaultName      string   `json:"vaultName"`
	VaultSlug      string   `json:"vaultSlug"`
}

func CreateVideos(videos []Video) ([]Video, error) {
	// We use make() here because at this point we know the size of
	// the slices and we won't need to reallocate memory if we were
	// to just loop and append.
	titles := make([]string, len(videos))
	slugs := make([]string, len(videos))
	studios := make([]string, len(videos))
	collectionIds := make([]int, len(videos))

	for i, v := range videos {
		titles[i] = v.Title
		slugs[i] = v.Slug
		studios[i] = v.Studio
		collectionIds[i] = v.CollectionID
	}

	query := `
		INSERT INTO videos (title, slug, studio, collection_id)
		SELECT *
		FROM UNNEST(
			$1::text[],
			$2::text[],
			$3::text[],
			$4::int[]
		)
		ON CONFLICT (slug, collection_id) DO UPDATE 
		SET
			title = EXCLUDED.title,
			studio = EXCLUDED.studio
		RETURNING id, title, slug, studio, collection_id
	`

	rows, err := db.Query(
		context.Background(),
		query,
		titles,
		slugs,
		studios,
		collectionIds,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Since we know the original size of the slice prior
	// to inserting, we know the max capacity of rows returned.
	// We don't specify length as a row conflict will result in
	// an update and not an insert.

	dbVideos := make([]Video, 0, len(videos))

	for rows.Next() {
		var v Video

		if err := rows.Scan(&v.ID, &v.Title, &v.Slug, &v.Studio, &v.CollectionID); err != nil {
			return nil, err
		}

		dbVideos = append(dbVideos, v)
	}

	return dbVideos, nil
}

func GetVideos(collectionId int) ([]Video, error) {
	query := `
		SELECT 
			v.id,
			v.title,
			v.slug,
			v.studio,
			c.name AS collection_name,
			c.slug AS collection_slug,
			va.id AS vault_id,
			va.slug AS vault_slug
		FROM 
			videos v
		JOIN 
			collections c ON v.collection_id = c.id
		JOIN 
			vaults va ON c.vault_id = va.id
		WHERE 
			c.id = $1
	`

	rows, err := db.Query(
		context.Background(),
		query,
		collectionId,
	)

	if err != nil {
		return nil, fmt.Errorf("fetching videos: %w", err)
	}
	defer rows.Close()

	var videos []Video

	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Slug, &v.Studio, &v.CollectionName, &v.CollectionSlug, &v.VaultID, &v.VaultSlug); err != nil {
			return nil, fmt.Errorf("fetching videos: %w", err)
		}

		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetching videos: %w", err)
	}

	return videos, nil
}

func GetVideo(videoId int) (*Video, error) {
	query := `
        SELECT 
            v.title,
            v.slug,
			v.studio,
            c.slug AS collection_slug,
            va.slug AS vault_slug,
			COALESCE(ARRAY_AGG(DISTINCT t.name) FILTER (WHERE t.name IS NOT NULL), '{}') AS tags,
			COALESCE(
				json_agg(
					DISTINCT jsonb_build_object('name', a.name, 'slug', a.slug)
				) FILTER (WHERE a.name IS NOT NULL),
				'[]'
			) AS actors
        FROM 
            videos v
        JOIN 
            collections c ON v.collection_id = c.id
        JOIN 
            vaults va ON c.vault_id = va.id
		LEFT JOIN 
    		video_tags vt ON vt.video_id = v.id
		LEFT JOIN 
			tags t ON t.id = vt.tag_id
		LEFT JOIN 
        	video_actors va2 ON va2.video_id = v.id
		LEFT JOIN 
        	actors a ON a.id = va2.actor_id
        WHERE 
            v.id = $1
		GROUP BY 
    		v.id, c.slug, va.slug
        LIMIT 1
    `

	var v Video

	err := db.QueryRow(
		context.Background(),
		query,
		videoId,
	).Scan(&v.Title, &v.Slug, &v.Studio, &v.CollectionSlug, &v.VaultSlug, &v.Tags, &v.Actors)

	if err != nil {
		return nil, fmt.Errorf("fetching video: %w", err)
	}

	return &v, nil
}
