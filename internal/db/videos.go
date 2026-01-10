package db

import (
	"context"
	"fmt"
	"log"

	"reelix-go/internal/utils"
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
	VaultID        int      `json:"vaultId"`
	VaultName      string   `json:"vaultName"`
}

func CreateVideo(video Video) error {
	tx, err := db.Begin(context.Background())

	if err != nil {
		return fmt.Errorf("failed to begin video transaction: %w", err)
	}

	defer tx.Rollback(context.Background())

	query := `
		INSERT INTO videos (title, slug, studio, collection_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (slug) DO UPDATE
		SET
			title = EXCLUDED.title,
			studio = EXCLUDED.studio
		RETURNING id
	`

	var videoId int

	err = tx.QueryRow(
		context.Background(),
		query,
		video.Title,
		video.Slug,
		video.Studio,
		video.CollectionID,
	).Scan(&videoId)

	if err != nil {
		return fmt.Errorf("db insert error: %w", err)
	}

	log.Printf("tags: %v (video: %v)", video.Tags, video.Title)

	for _, tag := range video.Tags {
		tagId, err := CreateTag(tag, tx)

		if err != nil {
			return fmt.Errorf("failed to create tag %v: %w", tag, err)
		}

		err = LinkVideoTag(videoId, *tagId, tx)

		if err != nil {
			return fmt.Errorf("failed to link tag %v to video %v: %w", tag, videoId, err)
		}
	}

	for _, actor := range video.Actors {
		var actorId *int
		var err error

		actorId, err = GetActor(actor.Name)

		if err != nil {
			newActor := Actor{
				Name: actor.Name,
				Slug: utils.TitleToSnake(actor.Name),
			}

			actorId, err = CreateActor(newActor)

			if err != nil {
				return fmt.Errorf("failed to create actor %v: %w", actor.Name, err)
			}
		}

		err = LinkVideoActor(videoId, *actorId, tx)

		if err != nil {
			return fmt.Errorf("failed to link actor %v to video %v: %w", actor.Name, videoId, err)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("video added: %v (collection: %v)", video.Title, video.CollectionID)

	return nil
}

func GetVideos(collectionId int) ([]Video, error) {
	query := `
		SELECT 
			v.id,
			v.title,
			v.slug,
			v.studio,
			c.name AS collection_name,
			va.id AS vault_id,
			va.name AS vault_name
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
		log.Fatal("videos query failed")
		return nil, err
	}
	defer rows.Close()

	var videos []Video

	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Slug, &v.Studio, &v.CollectionName, &v.VaultID, &v.VaultName); err != nil {
			log.Fatal("videos scan failed")
			return nil, err
		}

		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("videos rows failed")
		return nil, err
	}

	return videos, nil
}

func GetVideo(videoId int) (*Video, error) {
	query := `
        SELECT 
            v.title,
            v.slug,
			v.studio,
            c.name AS collection_name,
            va.name AS vault_name,
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
    		v.id, c.name, va.name
        LIMIT 1
    `

	var v Video

	err := db.QueryRow(
		context.Background(),
		query,
		videoId,
	).Scan(&v.Title, &v.Slug, &v.Studio, &v.CollectionName, &v.VaultName, &v.Tags, &v.Actors)

	if err != nil {
		return nil, fmt.Errorf("error fetching video: %v", err)
	}

	return &v, nil
}
