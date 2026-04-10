# Complete API Endpoint Reference

## Base Configuration

- **Base URL:** `http://localhost:8080`
- **API Version:** `/api/v1`
- **WebSocket URL:** `ws://localhost:8080/ws`
- **Authentication:** JWT Bearer Token in Authorization header

## Authentication Endpoints

### POST /api/v1/register - Register User

Register a new user account.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2026-04-07T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### POST /api/v1/login - Login User

Login with email and password to receive JWT token.

**Request:**
```json
{
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2026-04-07T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### GET /api/v1/health - Health Check

Check if the server is running and healthy.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Server is healthy",
  "data": {
    "status": "ok"
  }
}
```

## Profile Endpoints

### GET /api/v1/profile - Get User Profile

Get authenticated user's profile including wallet balance.

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "balance": 1000.50,
    "created_at": "2026-04-07T10:30:00Z"
  }
}
```

### PUT /api/v1/profile - Update User Profile

Update user's name and email.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "name": "Jane Doe",
    "email": "jane@example.com",
    "balance": 1000.50,
    "created_at": "2026-04-07T10:30:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Profile update failed",
  "error": "User with this email already exists"
}
```

## Wallet Endpoints

### GET /api/v1/wallet - Get Wallet Balance

Get user's current wallet balance and details.

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Wallet retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "balance": 1000.50,
    "currency": "USD",
    "created_at": "2026-04-07T10:30:00Z",
    "updated_at": "2026-04-07T10:30:00Z"
  }
}
```

## Betting Endpoints

### POST /api/v1/bet - Place a Bet

Place a new bet on the current game round.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "amount": 50.00,
  "auto_cashout": 2.5
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Bet placed successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "round_id": 1,
    "amount": 50.00,
    "auto_cashout": 2.5,
    "cashout_multiplier": 0.0,
    "result": "pending",
    "payout": 0.0,
    "created_at": "2026-04-07T10:30:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bet placement failed",
  "error": "Insufficient wallet balance"
}
```

### POST /api/v1/cashout - Cashout a Bet

Cashout (stop) an active bet at the current multiplier.

**Headers:**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "bet_id": 1
}
```

**Query Parameters:**
```
?multiplier=1.85
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Cashout successful",
  "data": {
    "bet_id": 1,
    "cashout_multiplier": 1.85,
    "payout": 92.50
  }
}
```

### GET /api/v1/bets/history - Get Bet History

Get user's betting history with pagination.

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
?page=1&limit=10
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Bet history retrieved",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "round_id": 1,
      "amount": 50.00,
      "auto_cashout": 2.5,
      "cashout_multiplier": 1.85,
      "result": "win",
      "payout": 92.50,
      "created_at": "2026-04-07T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 150
  }
}
```

## History Endpoints

### GET /api/v1/game/history - Get Game History

Get global game round history.

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
?page=1&limit=20
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Game history retrieved",
  "data": [
    {
      "id": 1,
      "round_code": "ROUND-2026-04-07-001",
      "crash_multiplier": 3.25,
      "status": "crashed",
      "start_time": "2026-04-07T10:30:00Z",
      "created_at": "2026-04-07T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 500
  }
}
```

### GET /api/v1/transactions/history - Get Transaction History

Get user's transaction history.

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
```
?page=1&limit=20
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Transaction history retrieved",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "type": "bet",
      "amount": 50.00,
      "description": "Bet placed on round ROUND-2026-04-07-001",
      "created_at": "2026-04-07T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 300
  }
}
```

## WebSocket Endpoint

### GET /ws - WebSocket Connection

Establishes a WebSocket connection for real-time game events.

**Connection URL:** `ws://localhost:8080/ws`

### WebSocket Events

#### countdown
Server sends countdown to game start.
```json
{
  "type": "countdown",
  "payload": {
    "seconds_remaining": 5
  }
}
```

#### round_start
Round has started, multiplier now growing.
```json
{
  "type": "round_start",
  "payload": {
    "round_code": "ROUND-2026-04-07-001",
    "round_id": 1
  }
}
```

#### multiplier_update
Multiplier has increased.
```json
{
  "type": "multiplier_update",
  "payload": {
    "round_code": "ROUND-2026-04-07-001",
    "multiplier": 2.45
  }
}
```

#### crash
Round has crashed.
```json
{
  "type": "crash",
  "payload": {
    "round_code": "ROUND-2026-04-07-001",
    "crash_multiplier": 5.67
  }
}
```

#### round_end
Round has ended.
```json
{
  "type": "round_end",
  "payload": {
    "round_code": "ROUND-2026-04-07-001",
    "crash_multiplier": 5.67,
    "total_bets": 42,
    "winners": 15
  }
}
```

#### active_bets
Active bet count update.
```json
{
  "type": "active_bets",
  "payload": {
    "round_code": "ROUND-2026-04-07-001",
    "active_bets": 42,
    "total_bet": 5000.00
  }
}
```

#### recent_winners
Recent winner notification.
```json
{
  "type": "recent_winners",
  "payload": {
    "username": "john_doe",
    "amount": 500.50,
    "multiplier": 2.50
  }
}
```

## Error Handling

### Common Error Response Format

All errors follow this format:

```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

### HTTP Status Codes

| Code | Meaning | Scenario |
|------|---------|----------|
| 200 | OK | Request succeeded |
| 201 | Created | Resource created (register, place bet) |
| 400 | Bad Request | Invalid input or validation failed |
| 401 | Unauthorized | Invalid or missing JWT token |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Email already exists |
| 500 | Server Error | Internal server error |
| 503 | Unavailable | Service temporarily unavailable |

### Common Error Codes

| Code | Message | Cause |
|------|---------|-------|
| USER_NOT_FOUND | User not found | User ID doesn't exist |
| INVALID_CREDENTIALS | Invalid email or password | Wrong login credentials |
| USER_EXISTS | User with this email already exists | Email already registered |
| INVALID_TOKEN | Invalid or expired token | JWT token is invalid/expired |
| UNAUTHORIZED | Unauthorized access | Missing/invalid authorization |
| INSUFFICIENT_BALANCE | Insufficient wallet balance | Bet amount exceeds balance |
| INVALID_BET_AMOUNT | Invalid bet amount | Bet amount is invalid |
| BET_NOT_FOUND | Bet not found | Bet ID doesn't exist |
| ROUND_NOT_FOUND | Round not found | Round ID doesn't exist |
| ROUND_CRASHED | Round has already crashed | Cannot perform operation on crashed round |
| ALREADY_CASHED_OUT | Bet already cashed out | Bet was already cashed out |

## Testing with cURL

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "Password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "Password123"
  }'
```

### Get Profile (with token)
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer {token}"
```

### Update Profile
```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Place Bet
```bash
curl -X POST http://localhost:8080/api/v1/bet \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50.00,
    "auto_cashout": 2.5
  }'
```

### Cashout Bet
```bash
curl -X POST http://localhost:8080/api/v1/cashout \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"bet_id": 1}' \
  -G -d "multiplier=1.85"
```

### Get Wallet
```bash
curl -X GET http://localhost:8080/api/v1/wallet \
  -H "Authorization: Bearer {token}"
```

### Get Bet History
```bash
curl -X GET "http://localhost:8080/api/v1/bets/history?page=1&limit=10" \
  -H "Authorization: Bearer {token}"
```

### Get Game History
```bash
curl -X GET "http://localhost:8080/api/v1/game/history?page=1&limit=20" \
  -H "Authorization: Bearer {token}"
```

### Get Transaction History
```bash
curl -X GET "http://localhost:8080/api/v1/transactions/history?page=1&limit=20" \
  -H "Authorization: Bearer {token}"
```

## API Response Summary

| Endpoint | Method | Auth | Response |
|----------|--------|------|----------|
| /health | GET | ✗ | `{status: "ok"}` |
| /register | POST | ✗ | `{user, token}` |
| /login | POST | ✗ | `{user, token}` |
| /profile | GET | ✓ | `{profile}` |
| /profile | PUT | ✓ | `{profile}` |
| /wallet | GET | ✓ | `{wallet}` |
| /bet | POST | ✓ | `{bet}` |
| /cashout | POST | ✓ | `{cashout_data}` |
| /bets/history | GET | ✓ | `{bets[], pagination}` |
| /game/history | GET | ✓ | `{rounds[], pagination}` |
| /transactions/history | GET | ✓ | `{transactions[], pagination}` |
| /ws | GET | ✗ | WebSocket connection |

## Rate Limiting

Currently, no rate limiting is implemented. For production, consider adding:
- Request rate limiting (e.g., 100 requests/minute per IP)
- Concurrent connection limits
- WebSocket message rate limiting

## Security Notes

1. **JWT Secret:** Change `JWT_SECRET` in production
2. **HTTPS:** Always use HTTPS in production
3. **CORS:** Configure appropriate CORS origins
4. **Password:** Minimum 6 characters, bcrypt hashed
5. **Tokens:** 24-hour expiration by default
6. **WebSocket:** Consider adding authentication
