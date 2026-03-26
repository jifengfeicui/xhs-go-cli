package detail

import (
	"path/filepath"
	"testing"

	"xhs-go-cli/internal/db"
)

type fakeClient struct{}

func (f *fakeClient) Detail(feedID string, xsecToken string) ([]byte, error) {
	return []byte(`{"success":true,"data":{"data":{"note_id":"` + feedID + `"}}}`), nil
}

func TestFetchAndStore(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "detail.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer database.Close()

	_, err = database.Exec(`INSERT INTO search_results(query_id, feed_id, xsec_token, title, author, raw_json) VALUES(1, 'feed1', 'token1', '标题1', '作者1', '{}')`)
	if err != nil {
		t.Fatalf("seed search_results: %v", err)
	}

	service := &Service{db: database, client: &fakeClient{}}
	rows, err := service.ListPending(10)
	if err != nil {
		t.Fatalf("list pending: %v", err)
	}
	result, err := service.FetchAndStore(rows, 2)
	if err != nil {
		t.Fatalf("fetch and store: %v", err)
	}
	if len(result) != 1 || result[0].Status != "ok" {
		t.Fatalf("unexpected result: %#v", result)
	}

	var count int
	if err := database.QueryRow(`SELECT COUNT(1) FROM details`).Scan(&count); err != nil {
		t.Fatalf("count details: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 detail row, got %d", count)
	}
}
