package db

import (
	"context"
	"fmt"
)

type Collection struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Path       string `json:"path"`
	VaultID    int    `json:"vaultId"`
	VaultName  string `json:"vaultName"`
	VaultSlug  string `json:"vaultSlug"`
	VideoCount int    `json:"videoCount"`
}

func CreateCollections(collections []Collection) ([]Collection, error) {
	// We use make() here because at this point we know the size of
	// the slices and we won't need to reallocate memory if we were
	// to just loop and append.

	names := make([]string, len(collections))
	slugs := make([]string, len(collections))
	paths := make([]string, len(collections))
	vaultIds := make([]int, len(collections))

	for i, c := range collections {
		names[i] = c.Name
		slugs[i] = c.Slug
		paths[i] = c.Path
		vaultIds[i] = c.VaultID
	}

	query := `
		INSERT INTO collections (name, slug, path, vault_id)
		SELECT *
		FROM UNNEST(
			$1::text[],
			$2::text[],
			$3::text[],
			$4::int[]
		)
		ON CONFLICT (slug, vault_id) 
		DO UPDATE SET
			name = EXCLUDED.name,
			path = EXCLUDED.path
		RETURNING id, name, slug, path, vault_id
	`

	rows, err := db.Query(
		context.Background(),
		query,
		names,
		slugs,
		paths,
		vaultIds,
	)

	if err != nil {
		return nil, fmt.Errorf("creating collections: %w", err)
	}

	defer rows.Close()

	// Since we know the original size of the slice prior
	// to inserting, we know the max capacity of rows returned.
	// We don't specify length as a row conflict will result in
	// an update and not an insert.

	dbCollections := make([]Collection, 0, len(collections))

	for rows.Next() {
		var c Collection

		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Path, &c.VaultID); err != nil {
			return nil, fmt.Errorf("creating collections: %w", err)
		}

		dbCollections = append(dbCollections, c)
	}

	return dbCollections, nil
}

func GetCollections(vaultId int) ([]Collection, error) {
	query := `
		SELECT 
			c.id, 
			c.name AS collection_name, 
			c.slug AS collection_slug,
			v.name AS vault_name,
			v.slug AS vault_slug,
			(
				SELECT COUNT(*)
				FROM videos vid
				WHERE vid.collection_id = c.id
			) AS video_count
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
		vaultId,
	)

	if err != nil {
		return nil, fmt.Errorf("fetching collections: %w", err)
	}
	defer rows.Close()

	var collections []Collection

	for rows.Next() {
		var c Collection

		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.VaultName, &c.VaultSlug, &c.VideoCount); err != nil {
			return nil, fmt.Errorf("fetching collections: %w", err)
		}

		collections = append(collections, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetching collections: %w", err)
	}

	return collections, nil
}

func GetCollection(collectionId int) (*Collection, error) {
	query := `
		SELECT 
			c.id, 
			c.name AS collection_name, 
			c.slug AS collection_slug,
			v.id AS vault_id,
			v.name AS vault_name,
			v.slug AS vault_slug,
			(
				SELECT COUNT(*)
				FROM videos vid
				WHERE vid.collection_id = c.id
			) AS video_count
		FROM 
			collections c
		JOIN 
			vaults v ON c.vault_id = v.id
		WHERE 
			c.id = $1
	`

	var c Collection

	err := db.QueryRow(
		context.Background(),
		query,
		collectionId,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.VaultID, &c.VaultName, &c.VaultSlug, &c.VideoCount)

	if err != nil {
		return nil, fmt.Errorf("fetching collection: %w", err)
	}

	return &c, nil
}
