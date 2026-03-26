package source

import (
	"os"
	"path/filepath"
	"testing"

	"xhs-go-cli/internal/db"
)

func TestImportFromJSON(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "import.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer database.Close()

	repo := NewRepo(database)
	input := `{
	  "records": [
	    {
	      "fields": {
	        "来源名称": [{"text": "兰蔻LANCOME"}],
	        "监控关键词": [{"text": "上海,快闪,打卡"}],
	        "来源类型": "品牌官方号",
	        "优先级": "高",
	        "名单层级": "白名单"
	      }
	    }
	  ]
	}`
	inputPath := filepath.Join(t.TempDir(), "sources.json")
	if err := os.WriteFile(inputPath, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	count, err := ImportFromJSON(repo, inputPath)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("got %d imported records", count)
	}

	sources, err := repo.List(10)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(sources) != 1 || sources[0].Name != "兰蔻LANCOME" {
		t.Fatalf("unexpected sources: %#v", sources)
	}
}
