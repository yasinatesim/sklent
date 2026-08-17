package invoicemodels

import "time"

// GIBSession is the single standing GIB portal login this app holds in process memory — never persisted to disk.
type GIBSession struct {
	Token     string
	CreatedAt time.Time
}
