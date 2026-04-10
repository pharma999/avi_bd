package utils

import (
	"time"
)

// GetCurrentTime returns the current time
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// FormatTime formats a time to a string
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses a time string
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// GetUnixTimestamp returns current unix timestamp
func GetUnixTimestamp() int64 {
	return time.Now().Unix()
}
