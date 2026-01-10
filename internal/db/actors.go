package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Actor struct {
	ID   int    `json:"id"`
	Name string `xml:"name" json:"name"`
	Slug string `json:"slug"`
}

func CreateActor(actor Actor) (*int, error) {
	query := `
		INSERT INTO actors (name, slug) VALUES ($1, $2)
		ON CONFLICT (name, slug) DO UPDATE SET 
			name = EXCLUDED.name,
			slug = EXCLUDED.slug
		RETURNING id
	`

	var actorId int

	err := db.QueryRow(
		context.Background(),
		query,
		actor.Name,
		actor.Slug,
	).Scan(&actorId)

	if err != nil {
		return nil, fmt.Errorf("failed to insert actor %s: %w", actor.Name, err)
	}

	log.Printf("actor added: %v", actor.Name)

	return &actorId, nil
}

func LinkVideoActor(videoId int, actorId int, tx pgx.Tx) error {
	query := `
		INSERT INTO video_actors (video_id, actor_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := tx.Exec(
		context.Background(),
		query,
		videoId,
		actorId,
	)

	if err != nil {
		return fmt.Errorf("failed to get actor %v to video: %w", actorId, err)
	}

	return nil
}

func GetActors(vaultId int) ([]Actor, error) {
	query := `
		SELECT
			a.id,
			a.name,
			a.slug
		FROM actors a
		WHERE EXISTS (
			SELECT 1
			FROM video_actors va
			JOIN videos v ON v.id = va.video_id
			JOIN collections c ON c.id = v.collection_id
			WHERE va.actor_id = a.id
			AND c.vault_id = $1
		)
		ORDER BY a.name
	`

	rows, err := db.Query(
		context.Background(),
		query,
		vaultId,
	)

	if err != nil {
		log.Fatal("actors query failed")
		return nil, err
	}
	defer rows.Close()

	var actors []Actor

	for rows.Next() {
		var a Actor
		if err := rows.Scan(&a.ID, &a.Name, &a.Slug); err != nil {
			log.Fatal("actors scan failed")
			return nil, err
		}

		actors = append(actors, a)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("actors rows failed")
		return nil, err
	}

	return actors, nil
}

func GetActor(name string) (*int, error) {
	query := `
		SELECT 
			id
		FROM
			actors
		WHERE	
			name = $1
		LIMIT 1
	`

	var a Actor

	err := db.QueryRow(
		context.Background(),
		query,
		name,
	).Scan(&a.ID)

	if err != nil {
		return nil, fmt.Errorf("error fetching actor: %v", err)
	}

	return &a.ID, nil
}
