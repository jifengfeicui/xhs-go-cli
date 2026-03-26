package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) Detail(feedID string, xsecToken string) ([]byte, error) {
	payload := map[string]any{"feed_id": feedID, "xsec_token": xsecToken}
	body, _ := json.Marshal(payload)
	resp, err := c.HTTP.Post(c.BaseURL+"/api/v1/feeds/detail", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("detail status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
