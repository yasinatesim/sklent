package promotion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinSplitIDs_RoundTrip(t *testing.T) {
	ids := []string{"p1", "p2", "p3"}
	assert.Equal(t, ids, SplitIDs(JoinIDs(ids)))
}

func TestSplitIDs_Empty(t *testing.T) {
	assert.Nil(t, SplitIDs(""))
}
