package qualify

import (
	"path/filepath"
	"testing"

	"xhs-go-cli/internal/db"
)

func TestQualifyAndStore(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "qualify.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer database.Close()

	detailJSON := `{"data":{"data":{"title":"兰蔻快闪｜来打卡啦！","note_url":"https://www.xiaohongshu.com/note/abc","desc":"上海新天地活动，到店打卡可领小样。需要先预约后参与领取。"}}}`
	_, err = database.Exec(`INSERT INTO details(feed_id, xsec_token, detail_json, fetch_status) VALUES('feed1', 'token1', ?, 'ok')`, detailJSON)
	if err != nil {
		t.Fatalf("seed details: %v", err)
	}

	service := NewService(database)
	rows, err := service.ListDetails(10)
	if err != nil {
		t.Fatalf("list details: %v", err)
	}
	result, err := service.QualifyAndStore(rows)
	if err != nil {
		t.Fatalf("qualify and store: %v", err)
	}
	if len(result) != 1 || result[0].Status != "accepted" {
		t.Fatalf("unexpected result: %#v", result)
	}

	var count int
	if err := database.QueryRow(`SELECT COUNT(1) FROM qualifications`).Scan(&count); err != nil {
		t.Fatalf("count qualifications: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 qualification row, got %d", count)
	}
}
