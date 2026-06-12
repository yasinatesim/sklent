package llm

import "context"

const (
	LLM_PROVIDER_OPENROUTER = "openrouter"
	LLM_PROVIDER_OFFLINE    = "offline"

	LLM_SHAPE_PRODUCT_COPY = "product_copy"
)

// Provider is the pluggable LLM surface. Implementations live in subpackages.
type Provider interface {
	Name() string
	Complete(ctx context.Context, system, user string) (string, error)
}

type registry struct {
	providers map[string]Provider
}

var defaultRegistry = &registry{providers: map[string]Provider{}}

func Register(p Provider) { defaultRegistry.providers[p.Name()] = p }

func Get(name string) (Provider, bool) {
	p, ok := defaultRegistry.providers[name]
	return p, ok
}
