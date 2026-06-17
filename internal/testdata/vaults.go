package testdata

import "reelix-go/internal/db"

var Vaults = []db.Vault{
	{
		Name: "Vault 1",
		Slug: "vault_1",
	},
	{
		Name: "Vault 2",
		Slug: "vault_2",
	},
}

func SeedVaults() ([]db.Vault, error) {
	dbVaults, err := db.CreateVaults(Vaults)

	if err != nil {
		return nil, err
	}

	return dbVaults, nil
}
