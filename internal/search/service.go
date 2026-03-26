package search

import (
	"database/sql"
	"encoding/json"

	"xhs-go-cli/internal/mcp"
)

type Service struct {
	db     *sql.DB
	client *mcp.Client
}

type QueryRow struct {
	ID    int64
	Query string
}

type SearchResult struct {
	FeedID    string `json:"id"`
	XsecToken string `json:"xsec_token"`
	Title     string `json:"title"`
	Author    string `json:"author"`
}

func NewService(db *sql.DB, client *mcp.Client) *Service {
	return &Service{db: db, client: client}
}

func (s *Service) ListQueries(limit int) ([]QueryRow, error) {
	rows, err := s.db.Query(`SELECT id, query FROM generated_queries ORDER BY id ASC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []QueryRow
	for rows.Next() {
		var q QueryRow
		if err := rows.Scan(&q.ID, &q.Query); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func (s *Service) SaveGeneratedQuery(sourceID int64, query string, queryType string) error {
	_, err := s.db.Exec(`INSERT INTO generated_queries(source_id, query, query_type) VALUES(?, ?, ?)`, sourceID, query, queryType)
	return err
}

func (s *Service) SearchAndStore(queryID int64, query string, limit int) (int, error) {
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
		_, err := s.db.Exec(
			`INSERT INTO search_results(query_id, feed_id, xsec_token, title, author, raw_json) VALUES(?, ?, ?, ?, ?, ?)`,
			queryID, item.FeedID, item.XsecToken, item.Title, item.Author, string(rawJSON),
		)
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
