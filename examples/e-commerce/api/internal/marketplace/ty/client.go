package ty

import (
	"context"
	"errors"
	"sync"
)

// Client is a skeleton marketplace client mirroring hb. Publish is a stub in the demo.
type Client struct {
	supplierID string
	baseURL    string

	mu       sync.RWMutex
	catCache []Category
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewClient(supplierID, baseURL string) *Client {
	return &Client{supplierID: supplierID, baseURL: baseURL}
}

func (c *Client) Categories(_ context.Context) ([]Category, error) {
	c.mu.RLock()
	if c.catCache != nil {
		defer c.mu.RUnlock()
		return c.catCache, nil
	}
	c.mu.RUnlock()

	cats := []Category{
		{ID: 2001, Name: "Moda > Tişört"},
		{ID: 2002, Name: "Takı & Mücevher > Bileklik"},
	}
	c.mu.Lock()
	c.catCache = cats
	c.mu.Unlock()
	return cats, nil
}

func (c *Client) Publish(_ context.Context, _ any) error {
	return errors.New("ty: publish not wired in demo")
}
