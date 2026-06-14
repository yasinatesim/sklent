package rag

import (
	"hash/fnv"
	"math"
	"strings"
)

const embedDim = 64

// HashingEmbedder is a dependency-free demo embedder; swap for a real embeddings API in production.
type Embedder interface {
	Embed(text string) []float32
}

type HashingEmbedder struct{}

func (HashingEmbedder) Embed(text string) []float32 {
	vec := make([]float32, embedDim)
	for _, tok := range strings.Fields(strings.ToLower(text)) {
		h := fnv.New32a()
		_, _ = h.Write([]byte(tok))
		vec[h.Sum32()%embedDim] += 1
	}
	norm := float32(0)
	for _, v := range vec {
		norm += v * v
	}
	if norm == 0 {
		return vec
	}
	inv := float32(1) / float32(math.Sqrt(float64(norm)))
	for i := range vec {
		vec[i] *= inv
	}
	return vec
}
