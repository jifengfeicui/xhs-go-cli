package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Search(query string, limit int) ([]byte, error) {
	payload := map[string]any{"keyword": query, "page": 1, "page_size": limit, "sort": "general"}
	body, _ := json.Marshal(payload)
	resp, err := c.HTTP.Post(c.BaseURL+"/api/v1/feeds/search", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("search status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
