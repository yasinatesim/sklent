package invoice

import (
	"context"
	"sync"
	"time"

	"github.com/yasinatesim/vela-commerce/api/internal/invoice/models"
)

// InMemoryGIBSessionStore holds the single standing GIB login in process memory — a restart naturally forces a
// fresh login rather than leaving a live GIB credential sitting in the database.
type InMemoryGIBSessionStore struct {
	mu      sync.Mutex
	current *invoicemodels.GIBSession
}

func NewInMemoryGIBSessionStore() *InMemoryGIBSessionStore {
	return &InMemoryGIBSessionStore{}
}

func (s *InMemoryGIBSessionStore) GetCurrentGIBSession(_ context.Context) (*invoicemodels.GIBSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current, nil
}

func (s *InMemoryGIBSessionStore) SaveGIBSession(_ context.Context, token string) (invoicemodels.GIBSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := invoicemodels.GIBSession{Token: token, CreatedAt: time.Now().UTC()}
	s.current = &row
	return row, nil
}

func (s *InMemoryGIBSessionStore) ClearGIBSession(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = nil
	return nil
}
