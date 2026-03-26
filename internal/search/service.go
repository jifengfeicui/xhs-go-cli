package search

import (
	"context"
	"encoding/json"
	"fmt"

	"xhs-go-cli/internal/mcp"
	"xhs-go-cli/internal/model"
	"xhs-go-cli/internal/repository"
)

type Service struct {
	queryRepo  repository.QueryRepository
	resultRepo repository.SearchResultRepository
	client     *mcp.Client
}

type QueryRow struct {
	ID    uint
	Query string
}

type SearchResult struct {
	FeedID    string `json:"id"`
	XsecToken string `json:"xsec_token"`
	Title     string `json:"title"`
	Author    string `json:"author"`
}

func NewService(queryRepo repository.QueryRepository, resultRepo repository.SearchResultRepository, client *mcp.Client) *Service {
	return &Service{
		queryRepo:  queryRepo,
		resultRepo: resultRepo,
		client:     client,
	}
}

func (s *Service) ListQueries(ctx context.Context, limit int) ([]QueryRow, error) {
	queries, err := s.queryRepo.List(ctx, limit)
	if err != nil {
		return nil, err
	}
	result := make([]QueryRow, len(queries))
	for i, q := range queries {
		result[i] = QueryRow{ID: q.ID, Query: q.Query}
	}
	return result, nil
}

func (s *Service) ListPendingQueries(ctx context.Context, limit int) ([]QueryRow, error) {
	queries, err := s.queryRepo.ListPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	result := make([]QueryRow, len(queries))
	for i, q := range queries {
		result[i] = QueryRow{ID: q.ID, Query: q.Query}
	}
	return result, nil
}

func (s *Service) SaveGeneratedQuery(ctx context.Context, sourceID uint, query string, queryType string) error {
	return s.queryRepo.Create(ctx, &model.GeneratedQuery{
		SourceID:  sourceID,
		Query:     query,
		QueryType: queryType,
	})
}

func (s *Service) SearchAndStore(ctx context.Context, queryID uint, query string, saveLimit int) (int, error) {
	filters := mcp.FilterOption{
		PublishTime: "一周内",
		NoteType:    "图文",
		Location:    "同城",
	}
	raw, err := s.client.Search(query, filters)
	if err != nil {
		return 0, err
	}
	items, err := parseSearchResults(raw)
	if err != nil {
		return 0, err
	}
	if saveLimit > 0 && len(items) > saveLimit {
		items = items[:saveLimit]
	}
	for _, item := range items {
		rawJSON, _ := json.Marshal(item)
		err := s.resultRepo.Create(ctx, &model.SearchResult{
			QueryID:   queryID,
			FeedID:    item.FeedID,
			XsecToken: item.XsecToken,
			Title:     item.Title,
			Author:    item.Author,
			RawJSON:   string(rawJSON),
		})
		if err != nil {
			return 0, err
		}
	}
	return len(items), nil
}

func parseSearchResults(raw []byte) ([]SearchResult, error) {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	data, _ := payload["data"].(map[string]any)
	feedsRaw, _ := data["feeds"].([]any)
	results := make([]SearchResult, 0, len(feedsRaw))
	for _, item := range feedsRaw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		noteCard, _ := m["noteCard"].(map[string]any)
		user, _ := noteCard["user"].(map[string]any)
		results = append(results, SearchResult{
			FeedID:    asString(m["id"]),
			XsecToken: asString(m["xsecToken"]),
			Title:     asString(noteCard["displayTitle"]),
			Author:    asString(user["nickname"]),
		})
	}
	return results, nil
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
