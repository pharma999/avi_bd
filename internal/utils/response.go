package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response structure per specification
type Response struct {
	Status    string      `json:"status"` // "success" or "error"
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorInfo represents error details per specification
type ErrorInfo struct {
	Code             int               `json:"code"`
	Message          string            `json:"message"`
	Field            *string           `json:"field,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// SuccessResponse sends a success response matching specification format
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:    "success",
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponse sends a simple error response (backward compatible with handlers)
// This automatically determines error code from status and handles both message strings and error objects
func ErrorResponse(c *gin.Context, statusCode int, message string, errOrVals interface{}) {
	code := GetErrorCodeFromStatus(statusCode)
	validationErrors := []ValidationError{}

	// Try to parse as validation errors
	if errs, ok := errOrVals.([]ValidationError); ok {
		validationErrors = errs
	}

	errInfo := &ErrorInfo{
		Code:             code,
		Message:          message,
		ValidationErrors: validationErrors,
	}

	c.JSON(statusCode, Response{
		Status:    "error",
		Error:     errInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponseWithField sends an error response with a specific field
func ErrorResponseWithField(c *gin.Context, statusCode int, errorCode int, message string, field string) {
	errInfo := &ErrorInfo{
		Code:    errorCode,
		Message: message,
		Field:   &field,
	}

	c.JSON(statusCode, Response{
		Status:    "error",
		Error:     errInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// PaginatedResponse represents a paginated response per specification
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination *Pagination `json:"pagination"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// PaginatedSuccessResponse sends paginated data with proper format
func PaginatedSuccessResponse(c *gin.Context, statusCode int, items interface{}, page int, limit int, total int64) {
	pages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division

	paginatedData := PaginatedResponse{
		Items: items,
		Pagination: &Pagination{
			Total: int(total),
			Page:  page,
			Limit: limit,
			Pages: pages,
		},
	}

	c.JSON(statusCode, Response{
		Status:    "success",
		Data:      paginatedData,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// GetErrorCodeFromStatus returns appropriate error code based on HTTP status
func GetErrorCodeFromStatus(statusCode int) int {
	switch statusCode {
	case 400:
		return 400 // Bad Request
	case 401:
		return 1001 // Invalid credentials/Unauthorized
	case 403:
		return 403 // Forbidden
	case 404:
		return 404 // Not Found
	case 409:
		return 2001 // Conflict (email exists)
	case 422:
		return 422 // Unprocessable Entity
	case 429:
		return 429 // Too Many Requests
	case 500:
		return 5001 // Database error
	case 503:
		return 503 // Service Unavailable
	default:
		return statusCode
	}
}

// GetTimestamp returns current timestamp in ISO 8601 format (UTC)
func GetTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
