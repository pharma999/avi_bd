package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

// UintToUUID converts a uint ID to a UUID-like string
// In production, consider using actual UUIDs in the database
// For now, this returns a formatted string
func UintToUUID(id uint) string {
	// Simple conversion: pad the uint to make it look like a UUID
	// Format: 00000000-0000-0000-0000-<uint padded to 12 chars>
	str := fmt.Sprintf("00000000-0000-5000-8000-%012d", id)
	return str
}

// StringToUint converts a string to uint
func StringToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// UUIDStringToUint converts a UUID-like string back to uint
// Extracts the trailing 12 digits which represent the original uint
func UUIDStringToUint(uuidStr string) (uint, error) {
	// Pattern: 00000000-0000-5000-8000-<12 digits>
	re := regexp.MustCompile(`-([0-9]{12})$`)
	matches := re.FindStringSubmatch(uuidStr)
	
	if len(matches) < 2 {
		// Try plain uint string
		return StringToUint(uuidStr)
	}
	
	val, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// UintToString converts a uint to string
func UintToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

// UintToPointerString converts a uint pointer to string pointer
func UintToPointerString(id *uint) *string {
	if id == nil {
		return nil
	}
	str := UintToString(*id)
	return &str
}

// Float64Ptr creates a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// StringPtr creates a pointer to a string
func StringPtr(s string) *string {
	return &s
}
