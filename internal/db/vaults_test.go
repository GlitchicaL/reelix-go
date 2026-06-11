package db

import (
	"testing"
)

func SeedVaults(t *testing.T) ([]Vault, error) {
	vaults := []Vault{
		{
			Name: "Vault 1",
			Slug: "vault_1",
		},
		{
			Name: "Vault 2",
			Slug: "vault_2",
		},
	}

	dbVaults, err := CreateVaults(vaults)

	if err != nil {
		return nil, err
	}

	return dbVaults, nil
}

func TestCreateVaults(t *testing.T) {
	SetupConnection(t)

	if _, err := SeedVaults(t); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetVaults(t *testing.T) {
	SetupConnection(t)
	vaults, _ := SeedVaults(t)

	dbVaults, err := GetVaults()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(vaults) != len(dbVaults) {
		t.Errorf("lengths of vaults do not match")
	}
}

func TestGetVault(t *testing.T) {
	SetupConnection(t)
	vaults, _ := SeedVaults(t)

	vault, err := GetVault(vaults[0].ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if vaults[0].Name != vault.Name {
		t.Errorf("vault %s does not match vault %s", vaults[0].Name, vault.Name)
	}
}
