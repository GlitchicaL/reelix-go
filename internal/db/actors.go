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

func CreateActors(actors []Actor) ([]Actor, error) {
	// We use make() here because at this point we know the size of
	// the slices and we won't need to reallocate memory if we were
	// to just loop and append.
	names := make([]string, len(actors))
	slugs := make([]string, len(actors))

	for i, a := range actors {
		names[i] = a.Name
		slugs[i] = a.Slug
	}

	query := `
		INSERT INTO actors (name, slug) 
		SELECT *
		FROM UNNEST(
			$1::text[],
			$2::text[]
		)
		ON CONFLICT (name, slug) DO UPDATE
		SET 
			name = EXCLUDED.name
		RETURNING id, name, slug
	`

	rows, err := db.Query(
		context.Background(),
		query,
		names,
		slugs,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Since we know the original size of the slice prior
	// to inserting, we know the max capacity of rows returned.
	// We don't specify length as a row conflict will result in
	// an update and not an insert.

	dbActors := make([]Actor, 0, len(actors))

	for rows.Next() {
		var a Actor

		if err := rows.Scan(&a.ID, &a.Name, &a.Slug); err != nil {
			return nil, err
		}

		dbActors = append(dbActors, a)
	}

	return dbActors, nil
}

func CreateActor(actor Actor) (*int, error) {
	query := `
		INSERT INTO actors (name, slug) VALUES ($1, $2)
		ON CONFLICT (name, slug) DO NOTHING
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

func LinkVideoActors(videos []Video, actors []Actor) error {
	// We use make() here because at this point we know the size of
	// the slices and we won't need to reallocate memory if we were
	// to just loop and append.
	videoIds := make([]int, 0, len(videos))
	actorIds := make([]int, 0, len(videos))

	/*
		Since videos may have more than 1 actor, we'll have to flatten
		the videoIds array. For example, if Video [1] has Actor [5, 8]:

		videoIds = [1, 1]
		actorIds = [5, 8]

		Since actor names are treated as unique, we can map the name to
		it's database ID
	*/

	actorsMap := map[string]int{}

	for _, a := range actors {
		actorsMap[a.Name] = a.ID
	}

	for _, v := range videos {
		if v.ID == 0 {
			continue
		}

		for _, a := range v.Actors {
			id, ok := actorsMap[a.Name]

			if !ok {
				continue
			}

			videoIds = append(videoIds, v.ID)
			actorIds = append(actorIds, id)
		}
	}

	query := `
		INSERT INTO video_actors (video_id, actor_id)
		SELECT *
		FROM UNNEST(
			$1::int[],
			$2::int[]
		)
		ON CONFLICT (video_id, actor_id) DO NOTHING
	`

	rows, err := db.Query(
		context.Background(),
		query,
		videoIds,
		actorIds,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

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
