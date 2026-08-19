package testdata

import (
	"reelix-go/internal/db"
)

var Collections = []db.Collection{
	{
		Name:    "Collection 1",
		Slug:    "collection_1",
		Path:    "/reelix/vaults/vault_1/collections/collection_1",
		VaultID: 1,
	},
	{
		Name:    "Collection 2",
		Slug:    "collection_2",
		Path:    "/reelix/vaults/vault_2/collections/collection_2",
		VaultID: 2,
	},
}

func SeedCollections() ([]db.Collection, error) {
	_, err := SeedVaults()

	if err != nil {
		return nil, err
	}

	dbCollections, err := db.CreateCollections(Collections)

	if err != nil {
		return nil, err
	}

	return dbCollections, nil
}
