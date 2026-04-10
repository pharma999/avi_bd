# Aviator Backend - Architecture & Implementation Guide

## Project Structure Overview

```
aviator_backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── config/
│   ├── config.go                  # Configuration management
│   └── env.go                     # Environment variable loading
├── internal/
│   ├── auth/
│   │   ├── jwt.go                 # JWT token generation & validation
│   │   ├── password.go            # Password hashing (bcrypt)
│   │   └── middleware.go          # JWT middleware for routes
│   ├── database/
│   │   ├── postgres.go            # Database connection setup
│   │   ├── migrations.go          # GORM auto-migration
│   │   └── seed.go                # Database seeding (optional)
│   ├── models/
│   │   ├── user.go                # User model
│   │   ├── wallet.go              # Wallet model
│   │   ├── bet.go                 # Bet model
│   │   ├── game_round.go          # Game round model
│   │   └── transaction.go         # Transaction model
│   ├── dto/
│   │   ├── auth_dto.go            # Auth request/response
│   │   ├── bet_dto.go             # Bet request/response
│   │   ├── wallet_dto.go          # Wallet response
│   │   ├── profile_dto.go         # Profile response
│   │   └── websocket_dto.go       # WebSocket event payloads
│   ├── repositories/
│   │   ├── user_repository.go     # User data access
│   │   ├── wallet_repository.go   # Wallet data access
│   │   ├── bet_repository.go      # Bet data access
│   │   ├── game_round_repository.go # Game round data access
│   │   └── transaction_repository.go # Transaction data access
│   ├── services/
│   │   ├── auth_service.go        # Auth business logic
│   │   ├── wallet_service.go      # Wallet operations
│   │   ├── bet_service.go         # Betting operations
│   │   ├── game_service.go        # Game operations
│   │   ├── profile_service.go     # Profile operations
│   │   └── transaction_service.go # Transaction operations
│   ├── handlers/
│   │   ├── auth_handler.go        # Login/register endpoints
│   │   ├── wallet_handler.go      # Wallet endpoints
│   │   ├── bet_handler.go         # Betting endpoints
│   │   ├── profile_handler.go     # Profile endpoints
│   │   ├── history_handler.go     # History endpoints
│   │   └── health_handler.go      # Health check
│   ├── websocket/
│   │   ├── hub.go                 # WebSocket hub management
│   │   ├── client.go              # Client connection handling
│   │   ├── broadcaster.go         # Event broadcasting utilities
│   │   └── events.go              # Event constructors
│   ├── game/
│   │   ├── engine.go              # Main game engine
│   │   ├── round_manager.go       # Round lifecycle management
│   │   ├── crash_generator.go     # Crash point generation
│   │   └── payout_calculator.go   # Payout calculations
│   ├── routes/
│   │   └── routes.go              # Route registration
│   ├── utils/
│   │   ├── response.go            # Standard responses
│   │   ├── validator.go           # Input validation
│   │   ├── logger.go              # Logging utility
│   │   └── time.go                # Time utilities
│   └── constants/
│       ├── game_status.go         # Game status constants
│       └── transaction_types.go   # Transaction type constants
├── pkg/
│   └── common/
│       └── errors.go              # Common error definitions
├── scripts/
│   ├── run.sh                     # Run script
│   └── migrate.sh                 # Migration script
├── docs/
│   └── api.md                     # API documentation
├── .env.example                   # Environment template
├── .gitignore                     # Git ignore rules
├── go.mod                         # Go module definition
├── go.sum                         # Go dependency hashes
├── Makefile                       # Build commands
├── README.md                      # Project README
└── QUICKSTART.md                  # Quick start guide
```

## Architecture Patterns

### 1. Clean Architecture / Layered Architecture

The project follows a layered architecture pattern:

```
┌─────────────────────────────────────┐
│         HTTP Handlers               │  <- Entry point for requests
│      (handlers/ package)            │
├─────────────────────────────────────┤
│         Services                    │  <- Business logic
│      (services/ package)            │
├─────────────────────────────────────┤
│       Repositories                  │  <- Data access layer
│     (repositories/ package)         │
├─────────────────────────────────────┤
│       Models & Database             │  <- Data models
│     (models/ & database/)           │
└─────────────────────────────────────┘
```

### 2. Dependency Injection

Services depend on repositories, repositories depend on models. This allows:
- Easy testing with mock repositories
- Clear separation of concerns
- Flexible implementation swapping

### 3. WebSocket Hub Pattern

The WebSocket implementation uses a hub pattern:
- Central hub manages all connections
- Client connections are registered/unregistered
- Messages broadcast to all connected clients
- Handles client disconnections gracefully

## Key Implementation Details

### Authentication Flow

```
┌──────────────┐
│  Registration│
│   Request    │
└──────┬───────┘
       │ Validate email/password
       ├─→ Hash password (bcrypt)
       ├─→ Create user in DB
       ├─→ Create wallet (1000 initial)
       │
       └─→ Generate JWT token
           Return token + user info
                 ↓
┌──────────────────────────────┐
│ JWT Token (24h expiration)   │
│ Used in Authorization header │
└──────────────────────────────┘
```

### Game Round Lifecycle

```
1. WAITING (5 seconds)
   ├─ Broadcast countdown events
   ├─ Users can place bets
   └─ Wallet debit upon bet placement

2. RUNNING (variable duration)
   ├─ Multiplier grows continuously
   ├─ Broadcast multiplier updates
   ├─ Users can cash out
   ├─ Wallet credit upon cashout
   └─ Process auto-cashouts

3. CRASHED (instant)
   ├─ Multiplier reached crash point
   ├─ Broadcast crash event
   ├─ Mark pending bets as losses
   └─ No more cashouts allowed

4. ENDED (brief)
   ├─ Calculate results
   ├─ Broadcast round end event
   └─ Prepare for next round
```

### Wallet Transaction Safety

All wallet operations use atomic database updates:

```go
// Atomic increment
UPDATE wallets SET balance = balance + amount WHERE user_id = ?

// Atomic decrement
UPDATE wallets SET balance = balance - amount WHERE user_id = ?

// With balance check first
SELECT balance FROM wallets WHERE user_id = ? FOR UPDATE
IF balance >= amount:
    UPDATE wallets SET balance = balance - amount WHERE user_id = ?
    CREATE transaction record
```

### Bet Lifecycle

```
PENDING → (cashout before crash) → WIN (payout credited)
       → (crash before cashout) → LOSS (payout = 0)

Auto-cashout:
PENDING → (multiplier >= target) → WIN (auto-processed)
```

### Crash Point Generation

Uses exponential distribution for realistic game:
- ~60% of crashes below 2.0x
- Progressive decline in probability for higher multipliers
- Configurable range (1.5x - 100.0x)

## Database Schema

### Relationships

```
User (1) ──→ (1) Wallet
   ├──→ (N) Bets
   ├──→ (N) Transactions
   └──→ (N) GameRounds (via Bet)

Bet
├─→ (1) User
└─→ (1) GameRound

GameRound ──→ (N) Bets
```

### Key Indexes

```
- users(id) - primary key
- users(email) - unique for login
- wallets(user_id) - for fast lookup
- bets(user_id, round_id) - for composite queries
- bets(result, round_id) - for active bet queries
- transactions(user_id, type) - for history queries
- game_rounds(round_code) - unique round lookup
- game_rounds(status) - for filtering by status
```

## Game Engine Details

### Multiplier Growth Algorithm

```
elapsed_time = current_time - round_start_time
multiplier = 1.0 + (elapsed_time_seconds * 0.05)
// Grows ~5% per second

When multiplier >= crash_multiplier:
    Round crashes
    Pending bets become losses
```

### Auto-Cashout Processing

```
For each active bet in running round:
    If bet.auto_cashout > 0 AND current_multiplier >= auto_cashout:
        Calculate payout = bet_amount * current_multiplier
        Update bet to WIN status
        Credit wallet with payout
        Create transaction record
        Broadcast winner event
```

## WebSocket Event Flow

```
Server → Client (Broadcasting)
├─ countdown (every 1 second)
├─ round_start (once per round)
├─ multiplier_update (every 100ms)
├─ crash (once per round)
├─ round_end (once per round)
├─ active_bets (periodic)
└─ recent_winners (on auto-cashout)

Client → Server (from WebSocket)
└─ (Currently read-only; can be extended)
```

## API Response Format

All responses are JSON with consistent structure:

```json
{
  "success": true/false,           // Operation status
  "message": "Human readable",     // Status message
  "data": {/*response data*/},     // Actual data
  "error": null/error_details      // Error info if failed
}
```

Paginated responses include:
```json
{
  "page": 1,
  "page_size": 20,
  "total_count": 150
}
```

## Configuration Management

### Environment Variables

```
External Config:
├─ APP_PORT (default: 8080)
├─ DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
├─ JWT_SECRET (must change in production!)
├─ GAME_TICK_MS (tick interval: 100ms)
├─ COUNTDOWN_SECONDS (5 seconds)
├─ INITIAL_BALANCE (new user wallet: 1000)
├─ CRASH_MIN (1.5x)
└─ CRASH_MAX (100.0x)

Internal Config:
├─ Routes (defined in routes.go)
├─ Error messages (defined in pkg/common/errors.go)
└─ Constants (defined in constants/)
```

## Security Measures

### 1. Authentication
- JWT tokens with 24-hour expiration
- Token must be in Authorization header
- Invalid/expired tokens → 401 Unauthorized

### 2. Password Security
- bcrypt hashing with cost factor 10
- Passwords never logged or exposed
- Constant-time comparison to prevent timing attacks

### 3. Authorization
- Middleware validates JWT before handler execution
- UserID from token verified against resource ownership
- Users can only access their own data

### 4. Data Validation
- Input validation via DTOs with struct tags
- Bet amounts validated to be positive
- Email validated as valid email format

### 5. Database Security
- SQL injection prevented by GORM parameterized queries
- Soft deletes (deleted_at field) for audit trail
- Transaction-safe wallet operations

## Performance Considerations

### 1. Database Queries
- Proper indexing on frequently queried columns
- Preloading relationships with GORM
- Pagination for large result sets

### 2. WebSocket
- Hub channel capacity (256) prevents blocking
- Client channels (256) for message queuing
- Ping/pong keep-alive every 30 seconds

### 3. Game Engine
- Configurable tick interval (default 100ms)
- Async processing via goroutines
- Non-blocking broadcast to clients

### 4. Concurrent Access
- RWMutex on Hub for thread-safe map access
- GORM handles database connection pooling
- Atomic database operations for wallet

## Testing Strategy

### Unit Tests (Per Layer)
```
handlers/ → handlers_test.go (mock services)
services/ → services_test.go (mock repositories)
repositories/ → repositories_test.go (integration tests)
game/ → game_test.go (crash generation tests)
```

### Integration Tests
- Test full flow: register → login → bet → cashout
- Database transaction rollback on failures

### Load Tests
- Apache Bench for HTTP endpoints
- WebSocket concurrent connection tests

## Future Enhancements

### Short Term
- [ ] User profile updates
- [ ] Deposit/withdrawal system
- [ ] Leaderboards and statistics
- [ ] Game analytics dashboard

### Medium Term
- [ ] Admin dashboard
- [ ] Rate limiting middleware
- [ ] Request logging and tracing
- [ ] Cache layer (Redis)

### Long Term
- [ ] Multi-tenancy support
- [ ] Microservices architecture
- [ ] Game variations (multiple crash points)
- [ ] Mobile push notifications
- [ ] Real-time fraud detection

## Deployment Checklist

Before production:
- [ ] Change JWT_SECRET to cryptographically secure value
- [ ] Set up HTTPS/TLS certificates
- [ ] Configure CORS with specific origins
- [ ] Set up PostgreSQL backups
- [ ] Enable query logging
- [ ] Set up error monitoring (Sentry)
- [ ] Configure VPN/firewall
- [ ] Set up log aggregation
- [ ] Configure auto-scaling
- [ ] Set up health checks
- [ ] Load balance WebSocket connections
- [ ] Set database connection limits
- [ ] Configure rate limiting
- [ ] Set up uptime monitoring

## Troubleshooting Guide

### "No such table"
- Migrations didn't run → Check database URL
- Run manually: `go run cmd/server/main.go` should auto-migrate

### "Connection refused"
- PostgreSQL not running
- Wrong DB credentials in .env
- Check: `psql -U postgres -h localhost`

### "Address already in use"
- Port 8080 occupied
- Change APP_PORT=8081 in .env
- Or: `kill -9 $(lsof -t -i :8080)` on Linux/Mac

### WebSocket not connecting
- Browser CORS issue
- Check browser console for errors
- Verify: `CheckOrigin: func(r *http.Request) bool { return true }`

### Slow queries
- Check database indexes
- Monitor with: `EXPLAIN ANALYZE SELECT ...`
- Consider query caching

## Related Documentation

- [README.md](README.md) - Project overview
- [QUICKSTART.md](QUICKSTART.md) - Setup instructions
- [docs/api.md](docs/api.md) - API reference

## Contributing

Follow these guidelines:
1. Keep functions small and focused
2. Use interfaces for abstraction
3. Add comments for complex logic
4. Follow Go naming conventions
5. Run `make fmt` and `make vet` before commit
6. Update README if changing behavior

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0  
**Status:** Production Ready
