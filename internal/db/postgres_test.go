package db_test

import (
	"reelix-go/internal/testutil"
	"testing"
)

var DB_TEST_URL = "postgres://testuser:testpass@localhost:5433/reelix-integration-db"

func TestConnection(t *testing.T) {
	if err := testutil.SetupConnection(DB_TEST_URL); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	t.Cleanup(testutil.Clean)

	testutil.Clean()
}
