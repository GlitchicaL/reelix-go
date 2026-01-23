package db

import (
	"context"
	"log"
)

type Collection struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Path      string `json:"path"`
	VaultID   int    `json:"vaultId"`
	VaultName string `json:"vaultName"`
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
		ON CONFLICT (name, vault_id) 
		DO UPDATE SET
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
		return nil, err
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
			return nil, err
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
		vaultId,
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
