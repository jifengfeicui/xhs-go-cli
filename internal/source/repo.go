package source

import (
	"context"
	"encoding/json"
	"os"

	"xhs-go-cli/internal/model"
	"xhs-go-cli/internal/repository"
)

type Source struct {
	ID             uint
	Name           string
	Keywords       string
	SourceType     string
	Priority       int
	Level          string
	City           string
	Confidence     string
	Reason         string
	PreferenceType string
	Remark         string
	LastStatus     string
}

type Repo struct {
	repo repository.SourceRepository
}

func NewRepo(repo repository.SourceRepository) *Repo {
	return &Repo{repo: repo}
}

func (r *Repo) Insert(ctx context.Context, src Source) (uint, error) {
	modelSrc := &model.Source{
		Name:           src.Name,
		Keywords:       src.Keywords,
		SourceType:     src.SourceType,
		Priority:       src.Priority,
		Level:          src.Level,
		City:           src.City,
		Confidence:     src.Confidence,
		Reason:         src.Reason,
		PreferenceType: src.PreferenceType,
		Remark:         src.Remark,
		LastStatus:     src.LastStatus,
	}
	err := r.repo.Create(ctx, modelSrc)
	if err != nil {
		return 0, err
	}
	return modelSrc.ID, nil
}

func (r *Repo) List(ctx context.Context, limit int) ([]Source, error) {
	sources, err := r.repo.List(ctx, limit)
	if err != nil {
		return nil, err
	}
	result := make([]Source, len(sources))
	for i, src := range sources {
		result[i] = Source{
			ID:             src.ID,
			Name:           src.Name,
			Keywords:       src.Keywords,
			SourceType:     src.SourceType,
			Priority:       src.Priority,
			Level:          src.Level,
			City:           src.City,
			Confidence:     src.Confidence,
			Reason:         src.Reason,
			PreferenceType: src.PreferenceType,
			Remark:         src.Remark,
			LastStatus:     src.LastStatus,
		}
	}
	return result, nil
}

type importPayload struct {
	Records []importRecord `json:"records"`
}

type importRecord struct {
	Fields map[string]any `json:"fields"`
}

func ImportFromJSON(ctx context.Context, repo *Repo, path string) (int, error) {
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
		_, err := repo.Insert(ctx, Source{
			Name:           name,
			Keywords:       plainText(fields["监控关键词"]),
			SourceType:     plainText(fields["来源类型"]),
			Priority:       toInt(fields["优先级"]),
			Level:          plainText(fields["名单层级"]),
			City:           plainText(fields["城市"]),
			Confidence:     plainText(fields["可信度"]),
			Reason:         plainText(fields["值得监控原因"]),
			PreferenceType: plainText(fields["偏向类型"]),
			Remark:         plainText(fields["备注"]),
			LastStatus:     plainText(fields["最近状态"]),
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

func toInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		var n int
		for _, c := range x {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			}
		}
		return n
	default:
		return 0
	}
}
