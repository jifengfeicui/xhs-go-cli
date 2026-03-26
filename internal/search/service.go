package search

import (
	"context"
	"encoding/json"

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

func (s *Service) SaveGeneratedQuery(ctx context.Context, sourceID uint, query string, queryType string) error {
	return s.queryRepo.Create(ctx, &model.GeneratedQuery{
		SourceID:  sourceID,
		Query:     query,
		QueryType: queryType,
	})
}

func (s *Service) SearchAndStore(ctx context.Context, queryID uint, query string, limit int) (int, error) {
	raw, err := s.client.Search(query, limit)
	if err != nil {
		return 0, err
	}
	items, err := parseSearchResults(raw)
	if err != nil {
		return 0, err
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
		return nil, err
	}
	data, _ := payload["data"].(map[string]any)
	inner, _ := data["data"].(map[string]any)
	itemsRaw, _ := inner["items"].([]any)
	results := make([]SearchResult, 0, len(itemsRaw))
	for _, item := range itemsRaw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		results = append(results, SearchResult{
			FeedID:    asString(m["id"]),
			XsecToken: asString(m["xsec_token"]),
			Title:     asString(m["title"]),
			Author:    asString(m["author"]),
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
