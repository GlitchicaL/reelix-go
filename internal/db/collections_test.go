package db_test

import (
	"reelix-go/internal/db"
	"reelix-go/internal/testdata"
	"reelix-go/internal/testutil"
	"testing"
)

func TestCreateCollections(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)

	if _, err := testdata.SeedCollections(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	testutil.Clean()
}

func TestGetCollections(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	collections, _ := testdata.SeedCollections()

	dbCollections, err := db.GetCollections(collections[0].VaultID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(dbCollections) == 0 {
		t.Errorf("collection fetch failed")
	}

	if collections[0].ID != dbCollections[0].ID {
		t.Errorf("collection ids do not match")
	}

	if collections[0].Name != dbCollections[0].Name {
		t.Errorf("collection names do not match")
	}

	if collections[0].Slug != dbCollections[0].Slug {
		t.Errorf("collection slugs do not match")
	}

	testutil.Clean()
}

func TestGetCollection(t *testing.T) {
	testutil.SetupConnection(DB_TEST_URL)
	collections, _ := testdata.SeedCollections()

	dbCollection, err := db.GetCollection(collections[0].ID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if collections[0].ID != dbCollection.ID {
		t.Errorf("collection id do not match")
	}

	if collections[0].Name != dbCollection.Name {
		t.Errorf("collection name do not match")
	}

	if collections[0].Slug != dbCollection.Slug {
		t.Errorf("collection slug do not match")
	}

	testutil.Clean()
}
