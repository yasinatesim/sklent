package invoice

import "fmt"

const sequenceWidth = 6

// NextNumber renders the human-readable, per-year sequential document number that goes on the
// GIB e-Arsiv invoice: PREFIX + year + zero-padded sequence.
func NextNumber(prefix string, year int, issuedThisYear int64) string {
	return fmt.Sprintf("%s%d%0*d", prefix, year, sequenceWidth, issuedThisYear+1)
}
