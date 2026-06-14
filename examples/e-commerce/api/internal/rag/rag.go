package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/yasinatesim/vela-commerce/api/internal/llm"
)

type EnhanceInput struct {
	RawTitle string
	Category string
}

type EnhanceOutput struct {
	SeoTitle      string `json:"seoTitle"`
	DescriptionTr string `json:"descriptionTr"`
	DescriptionEn string `json:"descriptionEn"`
	Slug          string `json:"slug"`
}

type Service struct {
	provider llm.Provider
	embedder Embedder
}

func New(provider llm.Provider) *Service {
	if provider == nil {
		provider = llm.OfflineProvider{}
	}
	return &Service{provider: provider, embedder: HashingEmbedder{}}
}

// Enhance retrieves similar products and asks the LLM for copy; offline it returns a deterministic draft.
func (s *Service) Enhance(ctx context.Context, in EnhanceInput, similar []string) (EnhanceOutput, error) {
	context := strings.Join(similar, "; ")
	user := fmt.Sprintf("%s\nCategory: %s\nSimilar: %s", in.RawTitle, in.Category, context)
	body, err := s.provider.Complete(ctx, productCopySystemPrompt, user)
	if err != nil {
		return EnhanceOutput{}, err
	}
	return EnhanceOutput{
		SeoTitle:      strings.TrimSpace(in.RawTitle),
		DescriptionTr: body,
		DescriptionEn: body,
		Slug:          Slugify(in.RawTitle),
	}, nil
}

const productCopySystemPrompt = "You write concise, SEO-friendly e-commerce product copy. " +
	"Return one short paragraph. Do not invent specifications you were not given."
