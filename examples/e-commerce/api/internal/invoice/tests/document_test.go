package invoice_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/invoice"
)

func TestNextNumberIsSequentialAndPrefixed(t *testing.T) {
	assert.Equal(t, "VC2026000001", invoice.NextNumber("VC", 2026, 0))
	assert.Equal(t, "VC2026000002", invoice.NextNumber("VC", 2026, 1))
	assert.Equal(t, "VC2026000100", invoice.NextNumber("VC", 2026, 99))
}

func TestNextNumberRestartsEachYear(t *testing.T) {
	assert.Equal(t, "VC2026000001", invoice.NextNumber("VC", 2026, 0))
	assert.Equal(t, "VC2027000001", invoice.NextNumber("VC", 2027, 0))
}
