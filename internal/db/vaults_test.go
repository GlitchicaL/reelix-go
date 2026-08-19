package db_test

import (
	"reelix-go/internal/db"
	"reelix-go/internal/testdata"
	"reelix-go/internal/testutil"
	"testing"
)

func TestCreateVaults(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)

	if _, err := testdata.SeedVaults(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	testutil.Clean()
}

func TestGetVaults(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	vaults, _ := testdata.SeedVaults()

	dbVaults, err := db.GetVaults()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(vaults) != len(dbVaults) {
		t.Errorf("lengths of vaults do not match")
	}

	testutil.Clean()
}

func TestGetVault(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	vaults, _ := testdata.SeedVaults()

	vault, err := db.GetVault(vaults[0].ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if vaults[0].Name != vault.Name {
		t.Errorf("vault %s does not match vault %s", vaults[0].Name, vault.Name)
	}

	testutil.Clean()
}
