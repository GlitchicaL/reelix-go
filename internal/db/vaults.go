package db

import (
	"context"
	"fmt"
)

type Vault struct {
	ID   int    `json:"vaultId"`
	Name string `json:"vaultName"`
}

func CreateVaults(vaults []Vault) ([]Vault, error) {
	// We use make() here because at this point we know the size of
	// the vault and we won't need to reallocate memory if we were
	// to just loop and append.
	names := make([]string, len(vaults))

	for i, v := range vaults {
		names[i] = v.Name
	}

	query := `
        INSERT INTO vaults (name)
        SELECT UNNEST($1::text[])
        ON CONFLICT (name)
        DO UPDATE SET name = EXCLUDED.name
        RETURNING id, name;
	`

	rows, err := db.Query(
		context.Background(),
		query,
		names,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dbVaults []Vault

	for rows.Next() {
		var v Vault

		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			return nil, err
		}

		dbVaults = append(dbVaults, v)
	}

	return dbVaults, nil
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

func GetVault(vaultId int) (*Vault, error) {
	query := `SELECT id, name FROM vaults WHERE id = $1`

	var va Vault

	err := db.QueryRow(
		context.Background(),
		query,
		vaultId,
	).Scan(&va.ID, &va.Name)

	if err != nil {
		return nil, fmt.Errorf("error fetching video: %v", err)
	}

	return &va, nil
}
