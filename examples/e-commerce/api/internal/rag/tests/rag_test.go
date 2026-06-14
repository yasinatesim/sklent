package rag_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/rag"
)

func TestSlugify(t *testing.T) {
	assert.Equal(t, "klasik-beyaz-tisort", rag.Slugify("Klasik Beyaz Tişört"))
	assert.Equal(t, "urun-1", rag.Slugify("  Ürün #1!  "))
}

func TestOfflineEnhance_NoKeyStillWorks(t *testing.T) {
	svc := rag.New(nil)
	out, err := svc.Enhance(context.Background(), rag.EnhanceInput{RawTitle: "Klasik Beyaz Tişört", Category: "Tişört"}, []string{"Beyaz tişört", "Pamuklu üst"})
	require.NoError(t, err)
	assert.Equal(t, "klasik-beyaz-tisort", out.Slug)
	assert.NotEmpty(t, out.DescriptionTr)
}

func TestHashingEmbedder_Normalized(t *testing.T) {
	e := rag.HashingEmbedder{}
	v := e.Embed("klasik beyaz tisort")
	var norm float32
	for _, x := range v {
		norm += x * x
	}
	assert.InDelta(t, 1.0, norm, 0.001)
}
