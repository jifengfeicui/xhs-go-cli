package source

import (
	"encoding/json"
	"os"
)

type importPayload struct {
	Records []importRecord `json:"records"`
}

type importRecord struct {
	Fields map[string]any `json:"fields"`
}

func ImportFromJSON(repo *Repo, path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	var payload importPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return 0, err
	}
	count := 0
	for _, record := range payload.Records {
		fields := record.Fields
		name := plainText(fields["来源名称"])
		if name == "" {
			continue
		}
		_, err := repo.Insert(Source{
			Name:       name,
			Keywords:   plainText(fields["监控关键词"]),
			SourceType: plainText(fields["来源类型"]),
			Priority:   plainText(fields["优先级"]),
			Level:      plainText(fields["名单层级"]),
		})
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func plainText(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []any:
		result := ""
		for _, item := range x {
			if m, ok := item.(map[string]any); ok {
				if text, ok := m["text"].(string); ok {
					result += text
				}
			}
		}
		return result
	default:
		return ""
	}
}
