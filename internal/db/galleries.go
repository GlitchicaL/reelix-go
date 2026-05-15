package db

import (
	"context"
	"fmt"
)

type Gallery struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	ImageCount int    `json:"imageCount"`
	VaultID    int    `json:"vaultId"`
	VaultSlug  string `json:"vaultSlug"`
}

func CreateGallery(galleries []Gallery) ([]Gallery, error) {
	// We use make() here because at this point we know the size of
	// the slices and we won't need to reallocate memory if we were
	// to just loop and append.

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
		ON CONFLICT (slug, vault_id) 
		DO UPDATE SET
			title = EXCLUDED.title,
			image_count = EXCLUDED.image_count
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
		return nil, fmt.Errorf("creating gallery: %w", err)
	}

	defer rows.Close()

	// Since we know the original size of the slice prior
	// to inserting, we know the max capacity of rows returned.
	// We don't specify length as a row conflict will result in
	// an update and not an insert.

	dbGalleries := make([]Gallery, 0, len(galleries))

	for rows.Next() {
		var g Gallery

		if err := rows.Scan(&g.ID, &g.Title, &g.Slug, &g.ImageCount, &g.VaultID); err != nil {
			return nil, fmt.Errorf("creating gallery: %w", err)
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
			v.slug AS vault_slug
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
		return nil, fmt.Errorf("fetching galleries: %w", err)
	}
	defer rows.Close()

	var galleries []Gallery

	for rows.Next() {
		var g Gallery

		if err := rows.Scan(&g.ID, &g.Title, &g.Slug, &g.ImageCount, &g.VaultID, &g.VaultSlug); err != nil {
			return nil, fmt.Errorf("fetching galleries: %w", err)
		}

		galleries = append(galleries, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetching galleries: %w", err)
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
			v.slug AS vault_slug
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
	).Scan(&g.ID, &g.Title, &g.Slug, &g.ImageCount, &g.VaultID, &g.VaultSlug)

	if err != nil {
		return nil, fmt.Errorf("fetching gallery: %w", err)
	}

	return &g, nil
}
