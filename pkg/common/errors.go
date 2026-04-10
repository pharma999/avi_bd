package common

import "fmt"

// AppError represents an application error
type AppError struct {
	Code    string
	Message string
	Details interface{}
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common errors
var (
	ErrUserNotFound        = AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrInvalidCredentials  = AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrUserExists          = AppError{Code: "USER_EXISTS", Message: "User with this email already exists"}
	ErrInvalidToken        = AppError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrUnauthorized        = AppError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrInsufficientBalance = AppError{Code: "INSUFFICIENT_BALANCE", Message: "Insufficient wallet balance"}
	ErrInvalidBetAmount    = AppError{Code: "INVALID_BET_AMOUNT", Message: "Invalid bet amount"}
	ErrBetNotFound         = AppError{Code: "BET_NOT_FOUND", Message: "Bet not found"}
	ErrRoundNotFound       = AppError{Code: "ROUND_NOT_FOUND", Message: "Round not found"}
	ErrRoundCrashed        = AppError{Code: "ROUND_CRASHED", Message: "Round has already crashed"}
	ErrAlreadyCashedOut    = AppError{Code: "ALREADY_CASHED_OUT", Message: "Bet already cashed out"}
	ErrInternalServer      = AppError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error"}
	ErrDatabaseError       = AppError{Code: "DATABASE_ERROR", Message: "Database error occurred"}
)

// NewAppError creates a new app error
func NewAppError(code, message string, details interface{}) AppError {
	return AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
