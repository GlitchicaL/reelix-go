package db

import (
	"context"
	"fmt"
	"log"
)

type Gallery struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	ImageCount int    `json:"imageCount"`
	VaultID    int    `json:"vaultId"`
	VaultName  string `json:"vaultName"`
}

func CreateGallery(galleries []Gallery) ([]Gallery, error) {
	titles := make([]string, len(galleries))
	slugs := make([]string, len(galleries))
	imageCounts := make([]int, len(galleries))
	vaultIds := make([]int, len(galleries))

	for i, g := range galleries {
		titles[i] = g.Title
		slugs[i] = g.Slug
		imageCounts[i] = g.ImageCount
		vaultIds[i] = g.VaultID
	}

	query := `
		INSERT INTO galleries (title, slug, image_count, vault_id)
		SELECT *
		FROM UNNEST(
			$1::text[],
			$2::text[],
			$3::int[],
			$4::int[]
		)
		ON CONFLICT (title, slug) 
		DO UPDATE SET
			title = EXCLUDED.title,
			slug = EXCLUDED.slug,
			image_count = EXCLUDED.image_count,
			vault_id = EXCLUDED.vault_id
		RETURNING id, title, slug, image_count, vault_id
	`

	rows, err := db.Query(
		context.Background(),
		query,
		titles,
		slugs,
		imageCounts,
		vaultIds,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dbGalleries []Gallery

	for rows.Next() {
		var g Gallery

		if err := rows.Scan(&g.ID, &g.Title, &g.Slug, &g.ImageCount, &g.VaultID); err != nil {
			return nil, err
		}

		dbGalleries = append(dbGalleries, g)
	}

	return dbGalleries, nil
}

func GetGalleries(vaultId int) ([]Gallery, error) {
	query := `
		SELECT 
			g.id,
			g.title,
			g.slug,
			g.image_count,
			v.id AS vault_id,
			v.name AS vault_name
		FROM
			galleries g
		JOIN 
			vaults v ON g.vault_id = v.id
		WHERE	
			g.vault_id = $1
	`

	rows, err := db.Query(
		context.Background(),
		query,
		vaultId,
	)

	if err != nil {
		log.Fatal("galleries query failed")
		return nil, err
	}
	defer rows.Close()

	var galleries []Gallery

	for rows.Next() {
		var g Gallery

		if err := rows.Scan(&g.ID, &g.Title, &g.Slug, &g.ImageCount, &g.VaultID, &g.VaultName); err != nil {
			log.Fatal("galleries scan failed")
			return nil, err
		}

		galleries = append(galleries, g)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("galleries rows failed")
		return nil, err
	}

	return galleries, nil
}

func GetGallery(galleryId int) (*Gallery, error) {
	query := `
		SELECT 
			g.id,
			g.title,
			g.slug,
			g.image_count,
			v.id AS vault_id,
			v.name AS vault_name
		FROM
			galleries g
		JOIN 
			vaults v ON g.vault_id = v.id
		WHERE	
			g.id = $1
	`

	var g Gallery

	err := db.QueryRow(
		context.Background(),
		query,
		galleryId,
	).Scan(&g.ID, &g.Title, &g.Slug, &g.ImageCount, &g.VaultID, &g.VaultName)

	if err != nil {
		return nil, fmt.Errorf("error fetching gallery: %v", err)
	}

	return &g, nil
}
