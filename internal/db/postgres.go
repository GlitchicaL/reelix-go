package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type Vault struct {
	ID   int    `json:"vaultId"`
	Name string `json:"vaultName"`
}

type Collection struct {
	ID        int    `json:"collectionId"`
	Name      string `json:"collectionName"`
	Path      string `json:"path"`
	VaultID   int    `json:"vaultId"`
	VaultName string `json:"vaultName"`
}

type Video struct {
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

type Actor struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func Connect(dbURL string) (*pgxpool.Pool, error) {
	var err error
	db, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return db, nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func CreateVault(vault Vault) error {
	query := `
		INSERT INTO vaults (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`

	_, err := db.Exec(
		context.Background(),
		query,
		vault.Name,
	)

	return err
}

func CreateCollection(collection Collection) error {
	query := `
		INSERT INTO collections (name, path, vault_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (name, vault_id) DO NOTHING
	`

	_, err := db.Exec(
		context.Background(),
		query,
		collection.Name,
		collection.Path,
		collection.VaultID,
	)

	if err != nil {
		return fmt.Errorf("db insert error: %w", err)
	}

	log.Printf("collection added: %v (vault: %v)", collection.Name, collection.VaultID)

	return nil
}

func CreateVideo(video Video) error {
	tx, err := db.Begin(context.Background())

	if err != nil {
		return fmt.Errorf("failed to begin video transaction: %w", err)
	}

	defer tx.Rollback(context.Background())

	query := `
		INSERT INTO videos (title, slug, studio, collection_id, vault_id)
		VALUES ($1, $2, $3, $4, $5)
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
		video.VaultID,
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

		err = LinkVideoTag(*tagId, videoId, tx)

		if err != nil {
			return fmt.Errorf("failed to link tag %v to video %v: %w", tag, videoId, err)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("video added: %v (collection: %v)", video.Title, video.CollectionID)

	return nil
}

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

func LinkVideoTag(tagId int, videoId int, tx pgx.Tx) error {
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

func GetVaults() ([]Vault, error) {
	query := `SELECT id, name FROM vaults`
	rows, err := db.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vaults []Vault

	for rows.Next() {
		var v Vault
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			return nil, err
		}
		vaults = append(vaults, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vaults, nil
}

func GetCollections(vaultId int) ([]Collection, error) {
	query := `
		SELECT 
			c.id, 
			c.name AS collection_name, 
			v.name AS vault_name
		FROM 
			collections c
		JOIN 
			vaults v ON c.vault_id = v.id
		WHERE 
			c.vault_id = $1
	`

	rows, err := db.Query(
		context.Background(),
		query,
		1,
	)

	if err != nil {
		log.Fatal("collection query failed")
		return nil, err
	}
	defer rows.Close()

	var collections []Collection

	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.VaultName); err != nil {
			return nil, err
		}

		log.Printf("collection: %v", c)

		collections = append(collections, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}

func GetVideos(collectionId int) ([]Video, error) {
	query := `
		SELECT 
			v.title,
			v.slug,
			v.studio,
			c.name AS collection_name,
			va.name AS vault_name
		FROM 
			videos v
		JOIN 
			collections c ON v.collection_id = c.id
		JOIN 
			vaults va ON v.vault_id = va.id
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
		if err := rows.Scan(&v.Title, &v.Slug, &v.Studio, &v.CollectionName, &v.VaultName); err != nil {
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

func GetVideo(collectionId int, videoSlug string) (*Video, error) {
	query := `
        SELECT 
            v.title,
            v.slug,
			v.studio,
            c.name AS collection_name,
            va.name AS vault_name,
			COALESCE(ARRAY_AGG(DISTINCT t.name) FILTER (WHERE t.name IS NOT NULL), '{}') AS tags
        FROM 
            videos v
        JOIN 
            collections c ON v.collection_id = c.id
        JOIN 
            vaults va ON v.vault_id = va.id
		LEFT JOIN 
    		video_tags vt ON vt.video_id = v.id
		LEFT JOIN 
			tags t ON t.id = vt.tag_id
        WHERE 
            v.collection_id = $1
        AND 
            v.slug = $2
		GROUP BY 
    		v.id, c.name, va.name
        LIMIT 1
    `

	var v Video

	err := db.QueryRow(
		context.Background(),
		query,
		collectionId,
		videoSlug,
	).Scan(&v.Title, &v.Slug, &v.Studio, &v.CollectionName, &v.VaultName, &v.Tags)

	if err != nil {
		return nil, fmt.Errorf("error fetching video: %v", err)
	}

	return &v, nil
}
