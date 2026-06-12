package llm

import (
	"context"
	"fmt"
	"strings"
)

// OfflineProvider returns a deterministic draft so the demo runs with no API key.
type OfflineProvider struct{}

func (OfflineProvider) Name() string { return LLM_PROVIDER_OFFLINE }

func (OfflineProvider) Complete(_ context.Context, _ string, user string) (string, error) {
	first := user
	if i := strings.IndexByte(user, '\n'); i > 0 {
		first = user[:i]
	}
	first = strings.TrimSpace(first)
	if first == "" {
		first = "Product"
	}
	return fmt.Sprintf("%s — quality and value, ready to ship.", first), nil
}
