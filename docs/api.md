# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <token>
```

## Response Format

All responses follow this format:

```json
{
  "success": true/false,
  "message": "Description",
  "data": { /* response data */ },
  "error": null/error_details
}
```

---

## Authentication Endpoints

### Register
**POST** `/register`

Register a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

**Errors:**
- `400`: Email already exists or validation failed

---

### Login
**POST** `/login`

Log in with credentials.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

**Errors:**
- `401`: Invalid email or password

---

## Profile Endpoints

### Get Profile
**GET** `/profile`

Get current user's profile information.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "balance": 1000.00,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Errors:**
- `401`: Unauthorized
- `404`: User not found

---

## Wallet Endpoints

### Get Wallet
**GET** `/wallet`

Get current wallet balance and info.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "message": "Wallet retrieved",
  "data": {
    "id": 1,
    "user_id": 1,
    "balance": 950.50,
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

**Errors:**
- `401`: Unauthorized
- `404`: Wallet not found

---

## Betting Endpoints

### Place Bet
**POST** `/bet`

Place a bet on the current game round.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "amount": 100.00,
  "auto_cashout": 2.50
}
```

- `amount`: Bet amount (required, > 0)
- `auto_cashout`: Optional multiplier to automatically cashout at

**Response (201):**
```json
{
  "success": true,
  "message": "Bet placed successfully",
  "data": {
    "id": 42,
    "user_id": 1,
    "round_id": 15,
    "amount": 100.00,
    "auto_cashout": 2.50,
    "cashout_multiplier": 0,
    "result": "pending",
    "payout": 0,
    "created_at": "2024-01-15T12:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

**Errors:**
- `400`: Invalid bet amount, insufficient balance, or waiting state not active
- `401`: Unauthorized
- `404`: Round not found

---

### Cashout Bet
**POST** `/cashout`

Cash out a pending bet at current multiplier.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Query Parameters:**
```
?multiplier=2.50
```

**Request Body:**
```json
{
  "bet_id": 42
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Cashout successful",
  "data": {
    "success": true,
    "message": "Cashout successful",
    "payout": 250.00,
    "new_balance": 1100.00
  }
}
```

**Errors:**
- `400`: Bet already cashed out, round has crashed, or other validation failure
- `401`: Unauthorized
- `403`: Bet doesn't belong to user
- `404`: Bet not found

---

### Get Bet History
**GET** `/bets/history`

Get user's betting history.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
```
?page=1&page_size=20
```

**Response (200):**
```json
{
  "success": true,
  "message": "Bet history retrieved",
  "data": [
    {
      "id": 42,
      "round_code": "ROUND_1234567890",
      "amount": 100.00,
      "result": "win",
      "payout": 250.00,
      "cashout_multiplier": 2.50,
      "crash_multiplier": 5.00,
      "created_at": "2024-01-15T12:00:00Z"
    },
    {
      "id": 41,
      "round_code": "ROUND_1234567889",
      "amount": 50.00,
      "result": "loss",
      "payout": 0,
      "cashout_multiplier": 0,
      "crash_multiplier": 1.50,
      "created_at": "2024-01-15T11:55:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total_count": 42
}
```

**Errors:**
- `401`: Unauthorized
- `500`: Server error

---

## History Endpoints

### Get Game History
**GET** `/game/history`

Get game round history.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
```
?page=1&page_size=20
```

**Response (200):**
```json
{
  "success": true,
  "message": "Game history retrieved",
  "data": [
    {
      "id": 15,
      "round_code": "ROUND_1234567890",
      "start_time": "2024-01-15T12:00:00Z",
      "crash_multiplier": 5.00,
      "status": "ended",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:15Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total_count": 150
}
```

**Errors:**
- `401`: Unauthorized
- `500`: Server error

---

### Get Transaction History
**GET** `/transactions/history`

Get user's transaction history.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
```
?page=1&page_size=50
```

**Response (200):**
```json
{
  "success": true,
  "message": "Transaction history retrieved",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "type": "bet",
      "amount": 100.00,
      "reference": "bet_42",
      "created_at": "2024-01-15T12:00:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "type": "cashout",
      "amount": 250.00,
      "reference": "bet_42",
      "created_at": "2024-01-15T12:00:10Z"
    }
  ],
  "page": 1,
  "page_size": 50,
  "total_count": 256
}
```

**Errors:**
- `401`: Unauthorized
- `500`: Server error

---

## Health Check

### Health Status
**GET** `/health`

Check server health.

**Response (200):**
```json
{
  "success": true,
  "message": "Server is healthy",
  "data": {
    "status": "ok"
  }
}
```

---

## WebSocket (Real-time Updates)

### Connect
**GET** `/ws`

Upgrade to WebSocket connection for real-time game updates.

**Example (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to game');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Event:', message.type, message.payload);
};

ws.onclose = () => {
  console.log('Disconnected from game');
};
```

---

## WebSocket Message Types

### countdown
Seconds remaining before round starts.
```json
{
  "type": "countdown",
  "payload": {
    "seconds_remaining": 5
  }
}
```

### round_start
Round has started.
```json
{
  "type": "round_start",
  "payload": {
    "round_code": "ROUND_1234567890",
    "round_id": 15
  }
}
```

### multiplier_update
Current multiplier (sent frequently).
```json
{
  "type": "multiplier_update",
  "payload": {
    "round_code": "ROUND_1234567890",
    "multiplier": 2.45
  }
}
```

### crash
Round has crashed at multiplier.
```json
{
  "type": "crash",
  "payload": {
    "round_code": "ROUND_1234567890",
    "crash_multiplier": 5.67
  }
}
```

### round_end
Round has ended.
```json
{
  "type": "round_end",
  "payload": {
    "round_code": "ROUND_1234567890",
    "crash_multiplier": 5.67,
    "total_bets": 42,
    "winners": 15
  }
}
```

### active_bets
Number of active bets in round.
```json
{
  "type": "active_bets",
  "payload": {
    "round_code": "ROUND_1234567890",
    "active_bets": 42,
    "total_bet": 5000.00
  }
}
```

### recent_winners
User just won.
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

---

## Error Codes

| Code | Status | Description |
|------|--------|-------------|
| 200  | OK | Request succeeded |
| 201  | Created | Resource created successfully |
| 400  | Bad Request | Invalid input or validation failed |
| 401  | Unauthorized | Missing or invalid token |
| 403  | Forbidden | Access denied |
| 404  | Not Found | Resource not found |
| 500  | Server Error | Internal server error |

---

## Rate Limiting

Currently no rate limiting is implemented. This should be added in production.

---

## Examples

### Complete Betting Flow

1. **Register User**
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice",
    "email": "alice@example.com",
    "password": "pass123"
  }'
```

2. **Login**
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "pass123"
  }'
```

3. **Get Wallet**
```bash
curl -X GET http://localhost:8080/api/v1/wallet \
  -H "Authorization: Bearer TOKEN"
```

4. **Place Bet**
```bash
curl -X POST http://localhost:8080/api/v1/bet \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "amount": 50.00,
    "auto_cashout": 2.00
  }'
```

5. **Cash Out** (must include current multiplier from WebSocket)
```bash
curl -X POST "http://localhost:8080/api/v1/cashout?multiplier=1.85" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "bet_id": 1
  }'
```

---

## Testing

Use Postman, Insomnia, or curl to test endpoints.

[Postman Collection Link](./postman_collection.json) - Import this file to Postman for easy API testing.
