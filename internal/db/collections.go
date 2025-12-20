package db

import (
	"context"
	"log"
)

type Collection struct {
	ID        int    `json:"collectionId"`
	Name      string `json:"collectionName"`
	Path      string `json:"path"`
	VaultID   int    `json:"vaultId"`
	VaultName string `json:"vaultName"`
}

func CreateCollections(collections []Collection) ([]Collection, error) {
	names := make([]string, len(collections))
	paths := make([]string, len(collections))
	vaultIds := make([]int, len(collections))

	for i, c := range collections {
		names[i] = c.Name
		paths[i] = c.Path
		vaultIds[i] = c.VaultID
	}

	query := `
		INSERT INTO collections (name, path, vault_id)
		SELECT *
		FROM UNNEST(
			$1::text[],
			$2::text[],
			$3::int[]
		)
		ON CONFLICT (name, vault_id) 
		DO UPDATE SET
			path = EXCLUDED.path
		RETURNING id, name, path, vault_id
	`

	rows, err := db.Query(
		context.Background(),
		query,
		names,
		paths,
		vaultIds,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dbCollections []Collection

	for rows.Next() {
		var c Collection

		if err := rows.Scan(&c.ID, &c.Name, &c.Path, &c.VaultID); err != nil {
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
