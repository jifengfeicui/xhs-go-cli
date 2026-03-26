package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenInitializesRequiredTables(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer database.Close()

	required := []string{"sources", "generated_queries", "search_results", "details", "qualifications"}
	for _, table := range required {
		if !tableExists(t, database, table) {
			t.Fatalf("expected table %s to exist", table)
		}
	}
}

func tableExists(t *testing.T, database *sql.DB, name string) bool {
	t.Helper()
	var count int
	if err := database.QueryRow(`SELECT COUNT(1) FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&count); err != nil {
		t.Fatalf("query sqlite_master failed: %v", err)
	}
	return count == 1
}
