package enhance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignificantTokens_DropsShort(t *testing.T) {
	got := SignificantTokens("Bir Klasik Beyaz Tişört")
	assert.Equal(t, []string{"bir", "klasik", "beyaz", "tişört"}, got)
}

func TestFilterCandidates_RanksByTokenHits(t *testing.T) {
	opts := []CategoryOption{
		{ID: 1, Name: "Giyim > Tişört"},
		{ID: 2, Name: "Aksesuar > Bileklik"},
		{ID: 3, Name: "Beyaz Eşya > Buzdolabı"},
	}
	got := FilterCandidates("Klasik Beyaz Tişört", opts, 5)
	assert.NotEmpty(t, got)
	assert.Equal(t, int64(1), got[0].ID)
}

func TestFilterCandidates_RespectsLimit(t *testing.T) {
	opts := []CategoryOption{
		{ID: 1, Name: "kirmizi tisort"},
		{ID: 2, Name: "kirmizi gomlek"},
		{ID: 3, Name: "kirmizi pantolon"},
	}
	got := FilterCandidates("kirmizi", opts, 2)
	assert.Len(t, got, 2)
}

func TestFilterCandidates_NoMatchReturnsNil(t *testing.T) {
	opts := []CategoryOption{{ID: 1, Name: "Mobilya"}}
	assert.Nil(t, FilterCandidates("Klasik Beyaz Tişört", opts, 5))
}
