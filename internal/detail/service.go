package detail

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"xhs-go-cli/internal/mcp"
	"xhs-go-cli/internal/model"
	"xhs-go-cli/internal/repository"
)

type detailFetcher interface {
	Detail(feedID string, xsecToken string) ([]byte, error)
}

type Service struct {
	resultRepo repository.SearchResultRepository
	detailRepo repository.DetailRepository
	client     detailFetcher
}

type SearchRow struct {
	ID        uint
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

func NewService(resultRepo repository.SearchResultRepository, detailRepo repository.DetailRepository, client *mcp.Client) *Service {
	return &Service{
		resultRepo: resultRepo,
		detailRepo: detailRepo,
		client:     client,
	}
}

func (s *Service) ListPending(ctx context.Context, limit int) ([]SearchRow, error) {
	results, err := s.resultRepo.ListPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	rows := make([]SearchRow, len(results))
	for i, r := range results {
		rows[i] = SearchRow{
			ID:        r.ID,
			FeedID:    r.FeedID,
			XsecToken: r.XsecToken,
			Title:     r.Title,
		}
	}
	return rows, nil
}

func (s *Service) FetchAndStore(ctx context.Context, rows []SearchRow, concurrency int) ([]FetchResult, error) {
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
			if err := s.saveDetail(ctx, row, raw, "ok"); err != nil {
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

func (s *Service) saveDetail(ctx context.Context, row SearchRow, raw []byte, status string) error {
	var compact map[string]any
	if err := json.Unmarshal(raw, &compact); err != nil {
		return fmt.Errorf("decode detail: %w", err)
	}
	body, _ := json.Marshal(compact)
	return s.detailRepo.Create(ctx, &model.Detail{
		FeedID:      row.FeedID,
		XsecToken:   row.XsecToken,
		DetailJSON:  string(body),
		FetchStatus: status,
	})
}
