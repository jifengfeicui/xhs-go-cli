package detail

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"xhs-go-cli/internal/mcp"
)


type detailFetcher interface {
	Detail(feedID string, xsecToken string) ([]byte, error)
}

type Service struct {
	db     *sql.DB
	client detailFetcher
}

type SearchRow struct {
	ID        int64
	FeedID    string
	XsecToken string
	Title     string
}

type FetchResult struct {
	FeedID string `json:"feed_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func NewService(db *sql.DB, client *mcp.Client) *Service {
	return &Service{db: db, client: client}
}

func (s *Service) ListPending(limit int) ([]SearchRow, error) {
	rows, err := s.db.Query(`SELECT id, feed_id, xsec_token, title FROM search_results ORDER BY id ASC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SearchRow
	for rows.Next() {
		var row SearchRow
		if err := rows.Scan(&row.ID, &row.FeedID, &row.XsecToken, &row.Title); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *Service) FetchAndStore(rows []SearchRow, concurrency int) ([]FetchResult, error) {
	if concurrency <= 0 {
		concurrency = 1
	}
	jobs := make(chan SearchRow)
	results := make(chan FetchResult, len(rows))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for row := range jobs {
			raw, err := s.client.Detail(row.FeedID, row.XsecToken)
			if err != nil {
				results <- FetchResult{FeedID: row.FeedID, Title: row.Title, Status: "error", Error: err.Error()}
				continue
			}
			if err := s.saveDetail(row, raw, "ok"); err != nil {
				results <- FetchResult{FeedID: row.FeedID, Title: row.Title, Status: "error", Error: err.Error()}
				continue
			}
			results <- FetchResult{FeedID: row.FeedID, Title: row.Title, Status: "ok"}
		}
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}
	for _, row := range rows {
		jobs <- row
	}
	close(jobs)
	wg.Wait()
	close(results)

	out := make([]FetchResult, 0, len(rows))
	for item := range results {
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) saveDetail(row SearchRow, raw []byte, status string) error {
	var compact map[string]any
	if err := json.Unmarshal(raw, &compact); err != nil {
		return fmt.Errorf("decode detail: %w", err)
	}
	body, _ := json.Marshal(compact)
	_, err := s.db.Exec(`INSERT INTO details(feed_id, xsec_token, detail_json, fetch_status) VALUES(?, ?, ?, ?)`, row.FeedID, row.XsecToken, string(body), status)
	return err
}
