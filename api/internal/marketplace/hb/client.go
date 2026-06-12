package hb

import (
	"context"
	"errors"
	"sync"
)

// Client is a skeleton marketplace client; categories are cached and Publish is a stub in the demo.
type Client struct {
	merchantID string
	baseURL    string

	mu       sync.RWMutex
	catCache []Category
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryAttribute struct {
	Name     string   `json:"name"`
	Required bool     `json:"required"`
	Enum     bool     `json:"enum"`
	Values   []string `json:"values,omitempty"`
}

func NewClient(merchantID, baseURL string) *Client {
	return &Client{merchantID: merchantID, baseURL: baseURL}
}

func (c *Client) Categories(_ context.Context) ([]Category, error) {
	c.mu.RLock()
	if c.catCache != nil {
		defer c.mu.RUnlock()
		return c.catCache, nil
	}
	c.mu.RUnlock()

	cats := []Category{
		{ID: 1001, Name: "Giyim > Tişört"},
		{ID: 1002, Name: "Aksesuar > Bileklik"},
		{ID: 1003, Name: "Ev > Mutfak"},
	}
	c.mu.Lock()
	c.catCache = cats
	c.mu.Unlock()
	return cats, nil
}

func (c *Client) CategoryAttributes(_ context.Context, categoryID int64) ([]CategoryAttribute, error) {
	if categoryID == 0 {
		return nil, errors.New("hb: categoryID required")
	}
	return []CategoryAttribute{
		{Name: "Renk", Required: true, Enum: true, Values: []string{"Beyaz", "Siyah", "Mavi"}},
		{Name: "Beden", Required: true, Enum: true, Values: []string{"S", "M", "L", "XL"}},
		{Name: "Materyal", Required: false, Enum: true, Values: []string{"Pamuk", "Polyester"}},
	}, nil
}

// Publish is intentionally a stub in the demo.
func (c *Client) Publish(_ context.Context, _ any) error {
	return errors.New("hb: publish not wired in demo")
}
