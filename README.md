# Aviator Backend

A production-ready backend for an Aviator-style realtime betting game built with Go, Gin, PostgreSQL, and WebSocket.

## Project Overview

This backend powers a realtime multiplayer betting game where users:
- Place bets on game rounds
- Watch multipliers increase in realtime
- Cash out before the round crashes
- Earn payouts based on their cash-out multiplier

The game features:
- **Countdown phase**: Users can place bets before round starts
- **Running phase**: Multiplier increases continuously, users can cash out
- **Crash phase**: Round crashes at server-determined crash point
- **Auto-cashout**: Optional automatic cashout when multiplier reaches target
- **Realtime updates**: WebSocket broadcasts game events to all connected clients

## Architecture

```
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── auth/            # JWT & password hashing
│   ├── database/        # Database connection & migrations
│   ├── models/          # Data models (User, Wallet, Bet, etc.)
│   ├── dto/             # Data Transfer Objects
│   ├── repositories/    # Database access layer
│   ├── services/        # Business logic layer
│   ├── handlers/        # HTTP request handlers
│   ├── websocket/       # WebSocket connection management
│   ├── game/            # Game engine (core Aviator logic)
│   ├── routes/          # Route definitions
│   ├── utils/           # Utility functions
│   └── constants/       # Application constants
├── pkg/common/          # Shared packages
└── scripts/             # Utility scripts
```

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL with GORM
- **Real-time**: WebSocket (Gorilla)
- **Authentication**: JWT (JSON Web Tokens)
- **Password**: bcrypt hashing

## Key Features

### 1. Authentication
- User registration and login
- JWT token generation and validation
- Password hashing with bcrypt
- Protected routes with middleware

### 2. Wallet Management
- User wallet balance tracking
- Transaction history recording
- Atomic balance updates (transaction-safe)
- Bet and cashout debit/credit operations

### 3. Betting System
- Place bets on active rounds
- Validate bet amounts and wallet balance
- Auto-cashout at specified multiplier
- Bet history with results and payouts

### 4. Game Engine
- Countdown phase before each round
- Continuous multiplier growth
- Server-side crash point determination
- Auto-cashout processing
- Realtime multiplier broadcasting
- Loss marking for uncashed bets

### 5. Real-time Updates
- WebSocket connections to all clients
- Live multiplier updates
- Countdown broadcasts
- Round start/crash/end events
- Active bet counts
- Recent winners display

## Database Schema

### Users Table
- `id` (Primary Key)
- `name`
- `email` (Unique)
- `password_hash`
- `created_at`, `updated_at`

### Wallets Table
- `id` (Primary Key)
- `user_id` (Foreign Key)
- `balance` (Decimal)
- `created_at`, `updated_at`

### Game Rounds Table
- `id` (Primary Key)
- `round_code` (Unique)
- `start_time`
- `crash_multiplier`
- `status` (waiting, running, crashed, ended)
- `created_at`, `updated_at`

### Bets Table
- `id` (Primary Key)
- `user_id` (Foreign Key)
- `round_id` (Foreign Key)
- `amount` (Decimal)
- `auto_cashout` (Decimal)
- `cashout_multiplier`
- `result` (pending, win, loss)
- `payout` (Decimal)
- `created_at`, `updated_at`

### Transactions Table
- `id` (Primary Key)
- `user_id` (Foreign Key)
- `type` (credit, debit, win, loss, etc.)
- `amount` (Decimal)
- `reference` (bet_id, round_id, etc.)
- `created_at`

## API Endpoints

### Authentication
```
POST   /api/v1/register          # Register new user
POST   /api/v1/login              # Login user
```

### Profile & Wallet
```
GET    /api/v1/profile            # Get user profile (Protected)
GET    /api/v1/wallet             # Get wallet balance (Protected)
```

### Betting
```
POST   /api/v1/bet                # Place a bet (Protected)
POST   /api/v1/cashout            # Cash out a bet (Protected)
GET    /api/v1/bets/history       # Get bet history (Protected)
```

### History
```
GET    /api/v1/game/history       # Get game round history (Protected)
GET    /api/v1/transactions/history # Get transaction history (Protected)
```

### Health
```
GET    /api/v1/health             # Health check
```

### WebSocket
```
GET    /ws                        # WebSocket connection
```

## WebSocket Events

### Events Sent by Server

**countdown**
```json
{
  "type": "countdown",
  "payload": {
    "seconds_remaining": 5
  }
}
```

**round_start**
```json
{
  "type": "round_start",
  "payload": {
    "round_code": "ROUND_1234567890",
    "round_id": 1
  }
}
```

**multiplier_update**
```json
{
  "type": "multiplier_update",
  "payload": {
    "round_code": "ROUND_1234567890",
    "multiplier": 2.45
  }
}
```

**crash**
```json
{
  "type": "crash",
  "payload": {
    "round_code": "ROUND_1234567890",
    "crash_multiplier": 5.67
  }
}
```

**round_end**
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

**active_bets**
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

**recent_winners**
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

## Request/Response Examples

### Register
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### Place Bet
```bash
curl -X POST http://localhost:8080/api/v1/bet \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "amount": 100.00,
    "auto_cashout": 2.50
  }'
```

### Get Wallet
```bash
curl -X GET http://localhost:8080/api/v1/wallet \
  -H "Authorization: Bearer TOKEN"
```

## Setup Instructions

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher
- git

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd avi_bd
```

2. **Create .env file**
```bash
cp .env.example .env
# Edit .env with your database credentials and config
```

3. **Install dependencies**
```bash
go mod download
go mod tidy
```

4. **Create PostgreSQL database**
```bash
psql -U postgres
CREATE DATABASE aviator_db;
\q
```

5. **Run the server**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Development

### Running Tests
```bash
go test ./...
go test ./... -v
```

### Running with hot reload (using air)
```bash
go install github.com/cosmtrek/air@latest
air
```

### Database Migrations
Migrations run automatically on server startup via GORM's `AutoMigrate()`.

To reset the database:
```bash
# Update main.go to uncomment database.DropAllTables(db)
# Then run the server
# Migrations will recreate all tables
```

## Environment Variables

```
APP_PORT=8080                              # Server port
DB_HOST=localhost                          # Database host
DB_PORT=5432                               # Database port
DB_USER=postgres                           # Database user
DB_PASSWORD=postgres                       # Database password
DB_NAME=aviator_db                         # Database name
JWT_SECRET=your-secret-key                 # JWT signing secret
GAME_TICK_MS=100                           # Game loop tick interval (ms)
COUNTDOWN_SECONDS=5                        # Countdown before round starts
INITIAL_BALANCE=1000.0                     # New user initial wallet balance
```

## Game Mechanics

### Round Lifecycle
1. **Waiting Phase** (5 seconds): Users place bets
2. **Running Phase**: Multiplier increases continuously
3. **Crash Point**: Server-determined crash occurs
4. **Ended Phase**: Round concludes, results calculated

### Crash Generation
- Uses weighted exponential distribution
- Lower multipliers more common (~60% under 2.0x)
- Configurable min/max range (1.5x - 100.0x)
- Ensures game fairness and realism

### Payout Calculation
- `Payout = Bet Amount × Cashout Multiplier`
- Results rounded to 2 decimal places
- Only paid if cashed out before crash
- Bet lost if round crashes before cashout

### Auto-Cashout
- If enabled and multiplier reaches target
- Automatically processes during running phase
- Funds credited to wallet immediately
- Notified via WebSocket event

## Security Considerations

1. **Password Security**: bcrypt with salting
2. **Authentication**: JWT tokens with expiration
3. **Database**: SQL injection protection via GORM
4. **CORS**: Configurable origin checking
5. **Rate Limiting**: Can be added via middleware
6. **Input Validation**: Request DTO validation
7. **Transaction Safety**: Atomic wallet operations

## Production Checklist

- [ ] Change `JWT_SECRET` to strong random value
- [ ] Enable HTTPS/TLS
- [ ] Set appropriate CORS origins
- [ ] Configure database backups
- [ ] Enable request logging
- [ ] Set up error monitoring (Sentry, etc.)
- [ ] Configure rate limiting
- [ ] Use connection pooling
- [ ] Set up database indexes
- [ ] Enable query caching
- [ ] Configure WebSocket timeouts
- [ ] Add authentication to WebSocket
- [ ] Set up health checks
- [ ] Configure auto-scaling
- [ ] Add request tracing
- [ ] Set up uptime monitoring

## Future Enhancements

- [ ] User profile updates
- [ ] Deposit/withdrawal system
- [ ] Leaderboards and statistics
- [ ] In-game chat
- [ ] Mobile app push notifications
- [ ] Admin dashboard
- [ ] Game analytics and reporting
- [ ] VIP levels and bonuses
- [ ] Referral system
- [ ] Multi-currency support
- [ ] Performance optimizations
- [ ] Load testing

## Troubleshooting

### Port already in use
```bash
# Change APP_PORT in .env or
lsof -i :8080  # Find process
kill -9 <pid>   # Kill process
```

### Database connection error
```bash
# Check PostgreSQL is running
psql -U postgres -h localhost -p 5432
# Verify credentials in .env
```

### WebSocket connection issues
```bash
# Check browser console for errors
# Verify WebSocket upgrade in network tab
# Check CORS settings
```

## License

MIT License - see LICENSE file for details

## Support

For issues and questions, please open an issue on the repository.
# avi_bd
