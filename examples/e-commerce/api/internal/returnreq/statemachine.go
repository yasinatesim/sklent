package returnreq

import "github.com/yasinatesim/vela-commerce/api/internal/constants"

// One table drives every question about the flow, so the allowed edges cannot drift apart from
// the UI's list of next actions.
var allowedTransitions = map[string][]string{
	constants.RETURN_STATUS_REQUESTED: {
		constants.RETURN_STATUS_APPROVED,
		constants.RETURN_STATUS_REJECTED,
	},
	constants.RETURN_STATUS_APPROVED: {
		constants.RETURN_STATUS_REFUNDED,
	},
	constants.RETURN_STATUS_REJECTED: {},
	constants.RETURN_STATUS_REFUNDED: {},
}

// A return may only be opened once the buyer actually has the goods.
var returnableOrderStatuses = map[string]struct{}{
	constants.ORDER_STATUS_SHIPPED: {},
}

func NextStatuses(from string) []string {
	return append([]string(nil), allowedTransitions[from]...)
}

func CanTransition(from, to string) bool {
	for _, candidate := range allowedTransitions[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

func IsTerminal(status string) bool {
	next, known := allowedTransitions[status]
	return known && len(next) == 0
}

func EligibleForReturn(orderStatus string) bool {
	_, ok := returnableOrderStatuses[orderStatus]
	return ok
}
