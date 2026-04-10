# API Implementation Guide - Aviator Betting Game Backend

## Overview
This document summarizes the work done to align the Go backend with the API specification document. The backend is a Gin-based REST API with WebSocket support for real-time game updates.

## Completed Tasks

### 1. ✅ DTOs Updated
All Data Transfer Objects have been updated to match the API specification exactly:

**Auth DTOs**
- `RegisterRequest`: Updated password validation (min 8, max 255)
- `LoginRequest`: Unchanged
- `AuthResponse`: Now includes `balance` and proper timestamp formatting
- `UserResponse`: Now uses string UUID format for ID and includes balance

**Profile DTOs**
- `ProfileResponse`: Uses string UUID for ID
- `UpdateProfileRequest`: Changed to optional fields (pointers) to support partial updates
- `StatisticsResponse`: Full statistics model per spec

**Wallet DTOs**
- `WalletResponse`: Uses string UUID for IDs, includes `currency` field

**Bet DTOs**
- `BetRequest`: Unchanged, validates amount > 0 and auto_cashout >= 1
- `BetResponse`: Added optional fields (cashout_multiplier, final_multiplier, winnings, ended_at)
- `CashoutRequest/CashoutResponse`: Full implementation per spec
- `BetHistoryItem/BetHistoryResponse`: Proper pagination support

**Transaction DTOs**
- `TransactionItem`: String UUIDs, all required fields per spec
- `TransactionHistoryResponse`: Paginated transactions list
- `TransactionRequest`: Updated UserID to string

**History DTOs**
- `GameRoundItem`: String UUID for ID
- `GameHistoryResponse`: Proper game history format
- `GameRoundResponse`: Full round state per spec
- `CurrentGameResponse`: Current game state for real-time updates

**WebSocket DTOs**
- Complete restructure matching specification events:
  - `AuthMessage`: Client authentication
  - `GameStartEvent`: sent to clients when game starts
  - `MultiplierUpdateEvent`: ~100ms updates during game
  - `GameCrashEvent`: Sent when game ends
  - `BetPlacedEvent`: Broadcast when bets are placed
  - `BetCashedOutEvent`: Broadcast when users cashout
  - `ErrorEvent`: Error handling
  - `PingMessage/PongMessage`: Keep-alive
  - `PlaceBetMessage/CashoutMessage`: Client messages

### 2. ✅ Utility Converters
Created `internal/utils/converter.go` with helper functions:
- `UintToUUID(id uint) string`: Converts database uint IDs to UUID-like strings
- `StringToUint(s string) uint`: Reverse conversion
- `UintToString(id uint) string`: Simple uint to string
- `Float64Ptr(f float64) *float64`: Helper for creating pointers
- `StringPtr(s string) *string`: Helper for string pointers

### 3. ✅ Services Updated

**AuthService** (`internal/services/auth_service.go`)
- Register: Now returns balance in response and uses UUID conversion
- Login: Returns balance and converted UUIDs
- VerifyToken: Returns user with UUID and balance
- All timestamps formatted as RFC3339

**ProfileService** (`internal/services/profile_service.go`)
- GetProfile: Uses UUID conversion, RFC3339 timestamps
- UpdateProfile: Handles optional fields, only updates provided values
- GetStatistics: Returns full statistics object
- Proper error handling with "email already exists" message

**WalletService** (`internal/services/wallet_service.go`)
- GetWallet: Returns UUID format IDs, includes currency, RFC3339 timestamps
- Credit/Debit: Transaction handling ready for wallet operations

### 4. ✅ Response Format
`internal/utils/response.go`:
- Standard Response wrapper with status, data, error, timestamp
- ErrorInfo with code, message, field, validation_errors
- PaginatedResponse structure wraps items and pagination metadata
- Error codes mapped correctly per specification

## Implementation Progress

### Remaining Service Updates Needed
The following services need similar UUID and timestamp updates (same pattern as AuthService):

1. **BetService** - Update PlaceBet, Cashout, GetBetHistory
2. **GameService** - Update game round responses  
3. **TransactionService** - Update transaction history responses

### Pattern for Remaining Services
Each service method should:
```go
// Example pattern
return &dto.SomeResponse{
    ID:        utils.UintToUUID(model.ID),
    UserID:    utils.UintToUUID(model.UserID),
    Timestamp: model.UpdatedAt.Format(time.RFC3339),
    // ... other fields
}
```

## Response Format Examples

### Success Response
```json
{
  "status": "success",
  "data": {
    "id": "00000000-0000-5000-8000-000000000001",
    "name": "John Doe",
    "email": "john@example.com",
    "balance": 1000.00,
    "created_at": "2024-04-07T10:30:00Z"
  },
  "timestamp": "2024-04-07T10:35:00Z"
}
```

### Paginated Response
```json
{
  "status": "success",
  "data": {
    "transactions": [
      {
        "id": "00000000-0000-5000-8000-000000000010",
        "user_id": "00000000-0000-5000-8000-000000000001",
        "type": "bet",
        "amount": 50.00,
        "balance_before": 1000.00,
        "balance_after": 950.00,
        "reference": "00000000-0000-5000-8000-000000000020",
        "description": "Placed bet on round",
        "status": "completed",
        "timestamp": "2024-04-07T10:30:00Z"
      }
    ],
    "pagination": {
      "total": 100,
      "page": 1,
      "limit": 10,
      "pages": 10
    }
  },
  "timestamp": "2024-04-07T10:35:00Z"
}
```

### Error Response
```json
{
  "status": "error",
  "error": {
    "code": 401,
    "message": "Invalid credentials",
    "validation_errors": []
  },
  "timestamp": "2024-04-07T10:35:00Z"
}
```

## WebSocket Implementation Guide

### Client Connection Flow
1. Client connects to `ws://localhost:8080/ws`
2. Client sends auth message: `{"type": "auth", "token": "jwt_token"}`
3. Server validates token and registers client
4. Client receives game events and can send commands

### Server Events
- `game_start`: New round begins
- `multiplier_update`: Current multiplier (every ~100ms)
- `game_crash`: Game ended with final multiplier
- `bet_placed`: User placed a bet
- `bet_cashed_out`: User cashed out
- `error`: Error occurred

### Client Commands
- `place_bet`: Submit bet
- `cashout`: Manually cashout
- `ping`: Keep connection alive

## API Endpoints Implemented

### Authentication
- ✅ POST `/api/v1/register` - Register user
- ✅ POST `/api/v1/login` - Login user
- ✅ GET `/api/v1/verify` - Verify token
- ✅ GET `/api/v1/health` - Health check

### Profile
- ✅ GET `/api/v1/profile` - Get user profile
- ✅ PUT `/api/v1/profile` - Update profile
- ✅ GET `/api/v1/profile/statistics` - Get statistics

### Wallet
- ✅ GET `/api/v1/wallet` - Get wallet

### Betting
- ✅ POST `/api/v1/bet` - Place bet
- ✅ POST `/api/v1/cashout` - Cashout bet
- ✅ GET `/api/v1/bets/history` - Bet history

### History
- ✅ GET `/api/v1/game/history` - Game history
- ✅ GET `/api/v1/transactions/history` - Transaction history

### WebSocket
- ✅ GET `/ws` - WebSocket connection

## Testing Recommendations

### 1. Authentication Flow
```bash
# Register
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "SecurePass123!"}'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "SecurePass123!"}'

# Verify token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/verify
```

### 2. Pagination
All list endpoints support pagination with defaults:
- `page`: Default 1
- `limit`: Default 10-50 (endpoint specific), max 100
- Returns `total`, `page`, `limit`, `pages` in pagination metadata

### 3. Error Handling
All endpoints return proper error codes:
- 400: Bad Request (validation failures)
- 401: Unauthorized (missing/invalid token)
- 404: Not Found (resource doesn't exist)
- 409: Conflict (email already exists)
- 422: Unprocessable Entity (validation failed)

## Configuration

### Environment Variables
- `DB_URL`: Database connection string
- `JWT_SECRET`: JWT signing secret
- `JWT_EXPIRY`: Token expiration time
- `INITIAL_BALANCE`: Starting wallet balance
- `PORT`: Server port (default 8080)

### Database
- Uses PostgreSQL with GORM
- All migrations run on startup
- Decimal(20,2) for monetary values
- UUID-like strings in API (stored as uint in DB)

## Next Steps

1. **Complete Service Updates**: Apply UUID conversion pattern to remaining services
2. **WebSocket Handler**: Implement WebSocket event broadcasting
3. **Integration Tests**: Test full workflow end-to-end
4. **Performance**: Add database indexes for common queries
5. **Documentation**: Generate Swagger docs from handlers

## Files Modified
- `internal/dto/auth_dto.go` ✅
- `internal/dto/profile_dto.go` ✅
- `internal/dto/wallet_dto.go` ✅
- `internal/dto/bet_dto.go` ✅
- `internal/dto/transaction_dto.go` ✅
- `internal/dto/history_dto.go` ✅
- `internal/dto/response_dto.go` ✅
- `internal/dto/websocket_dto.go` ✅
- `internal/services/auth_service.go` ✅
- `internal/services/profile_service.go` ✅
- `internal/services/wallet_service.go` ✅
- `internal/utils/converter.go` ✅ (created)
- `internal/utils/response.go` (no changes needed)

## Files Needing Updates
- `internal/services/bet_service.go` (high priority)
- `internal/services/game_service.go` (high priority)
- `internal/services/transaction_service.go` (high priority)
- Handlers may need adjustments based on service changes
- WebSocket implementation in `internal/websocket/`

---

**Last Updated**: 2024-04-07
**Status**: ~60% Complete - Core foundation in place, services integration in progress
