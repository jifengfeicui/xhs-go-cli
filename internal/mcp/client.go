package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FilterOption struct {
	SortBy      string `json:"sort_by,omitempty"`
	NoteType    string `json:"note_type,omitempty"`
	PublishTime string `json:"publish_time,omitempty"`
	SearchScope string `json:"search_scope,omitempty"`
	Location    string `json:"location,omitempty"`
}

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Search(query string, filters FilterOption) ([]byte, error) {
	payload := map[string]any{
		"keyword": query,
		"page":    1,
	}
	if filters.PublishTime != "" || filters.NoteType != "" || filters.Location != "" || filters.SearchScope != "" {
		payload["filters"] = filters
	}
	body, _ := json.Marshal(payload)
	resp, err := c.HTTP.Post(c.BaseURL+"/api/v1/feeds/search", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("search status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
