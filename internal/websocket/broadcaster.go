package websocket

import (
	"time"
)

// NewDeadline returns a new write deadline
func NewDeadline() time.Time {
	return time.Now().Add(10 * time.Second)
}

// NewTicker returns a new ticker for ping messages
func NewTicker() *time.Ticker {
	return time.NewTicker(30 * time.Second)
}
