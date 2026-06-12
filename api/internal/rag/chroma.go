package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// ChromaClient wraps ChromaDB's HTTP API so retrieval lives inside the Go API, no extra service.
type ChromaClient struct {
	baseURL string
	http    *http.Client
}

func NewChromaClient(baseURL string) *ChromaClient {
	return &ChromaClient{baseURL: baseURL, http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *ChromaClient) Heartbeat(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/heartbeat", nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *ChromaClient) post(ctx context.Context, path string, payload any) (*http.Response, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.http.Do(req)
}

var _ = (*ChromaClient)(nil).post
