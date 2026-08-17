package invoice_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/invoice"
)

func TestGIBSessionStoreStartsEmptySoARestartForcesAFreshLogin(t *testing.T) {
	store := invoice.NewInMemoryGIBSessionStore()

	current, err := store.GetCurrentGIBSession(context.Background())

	require.NoError(t, err)
	assert.Nil(t, current, "a fresh process must hold no GIB credential")
}

func TestGIBSessionStoreKeepsOnlyTheLatestLogin(t *testing.T) {
	store := invoice.NewInMemoryGIBSessionStore()
	ctx := context.Background()

	_, err := store.SaveGIBSession(ctx, "first")
	require.NoError(t, err)
	saved, err := store.SaveGIBSession(ctx, "second")
	require.NoError(t, err)
	assert.Equal(t, "second", saved.Token)

	current, err := store.GetCurrentGIBSession(ctx)
	require.NoError(t, err)
	require.NotNil(t, current)
	assert.Equal(t, "second", current.Token)
	assert.False(t, current.CreatedAt.IsZero())
}

func TestGIBSessionStoreClears(t *testing.T) {
	store := invoice.NewInMemoryGIBSessionStore()
	ctx := context.Background()
	_, err := store.SaveGIBSession(ctx, "tok")
	require.NoError(t, err)

	require.NoError(t, store.ClearGIBSession(ctx))

	current, err := store.GetCurrentGIBSession(ctx)
	require.NoError(t, err)
	assert.Nil(t, current)
}

func TestGIBSessionStoreIsSafeUnderConcurrentUse(t *testing.T) {
	store := invoice.NewInMemoryGIBSessionStore()
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); _, _ = store.SaveGIBSession(ctx, "tok") }()
		go func() { defer wg.Done(); _, _ = store.GetCurrentGIBSession(ctx) }()
	}
	wg.Wait()

	current, err := store.GetCurrentGIBSession(ctx)
	require.NoError(t, err)
	require.NotNil(t, current)
}
