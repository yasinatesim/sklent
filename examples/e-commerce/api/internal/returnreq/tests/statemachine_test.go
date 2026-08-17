package returnreq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/returnreq"
)

func TestCanTransition(t *testing.T) {
	cases := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{"requested to approved", constants.RETURN_STATUS_REQUESTED, constants.RETURN_STATUS_APPROVED, true},
		{"requested to rejected", constants.RETURN_STATUS_REQUESTED, constants.RETURN_STATUS_REJECTED, true},
		{"approved to refunded", constants.RETURN_STATUS_APPROVED, constants.RETURN_STATUS_REFUNDED, true},
		{"rejected is terminal", constants.RETURN_STATUS_REJECTED, constants.RETURN_STATUS_APPROVED, false},
		{"refunded is terminal", constants.RETURN_STATUS_REFUNDED, constants.RETURN_STATUS_APPROVED, false},
		{"cannot skip approval", constants.RETURN_STATUS_REQUESTED, constants.RETURN_STATUS_REFUNDED, false},
		{"cannot go backwards", constants.RETURN_STATUS_APPROVED, constants.RETURN_STATUS_REQUESTED, false},
		{"same status is not a transition", constants.RETURN_STATUS_APPROVED, constants.RETURN_STATUS_APPROVED, false},
		{"unknown source", "bogus", constants.RETURN_STATUS_APPROVED, false},
		{"unknown target", constants.RETURN_STATUS_REQUESTED, "bogus", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, returnreq.CanTransition(tc.from, tc.to))
		})
	}
}

func TestIsTerminal(t *testing.T) {
	assert.False(t, returnreq.IsTerminal(constants.RETURN_STATUS_REQUESTED))
	assert.False(t, returnreq.IsTerminal(constants.RETURN_STATUS_APPROVED))
	assert.True(t, returnreq.IsTerminal(constants.RETURN_STATUS_REJECTED))
	assert.True(t, returnreq.IsTerminal(constants.RETURN_STATUS_REFUNDED))
}

func TestEligibleForReturn(t *testing.T) {
	cases := []struct {
		name        string
		orderStatus string
		want        bool
	}{
		{"shipped orders can be returned", constants.ORDER_STATUS_SHIPPED, true},
		{"paid but unshipped cannot", constants.ORDER_STATUS_PAID, false},
		{"pending cannot", constants.ORDER_STATUS_PENDING, false},
		{"cancelled cannot", constants.ORDER_STATUS_CANCELLED, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, returnreq.EligibleForReturn(tc.orderStatus))
		})
	}
}

func TestNextStatusesIsDerivedFromTheSameTable(t *testing.T) {
	next := returnreq.NextStatuses(constants.RETURN_STATUS_REQUESTED)
	require.Len(t, next, 2)
	assert.ElementsMatch(t,
		[]string{constants.RETURN_STATUS_APPROVED, constants.RETURN_STATUS_REJECTED}, next)

	assert.Empty(t, returnreq.NextStatuses(constants.RETURN_STATUS_REFUNDED))
}
