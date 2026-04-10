# Backend Management Quick Reference
## Aviator Betting Game - Go Backend

**Location**: `D:\go project\avi_bd`  
**Status**: 70% Complete  
**Framework**: Go + Gin + PostgreSQL  
**Last Updated**: 2024-04-07

---

## 🚀 Quick Start (Copy & Paste)

### First Time Setup
```bash
cd D:\go project\avi_bd

# Install dependencies
go mod download

# Create database
psql -U postgres -c "CREATE DATABASE aviator_db;"

# Run server
make run
# OR: go run cmd/server/main.go
```

### Daily Development
```bash
# Start server (keep running in one terminal)
make run

# In another terminal, test API
curl http://localhost:8080/api/v1/health

# Make code changes, auto-reload should happen
# Then use Postman or cURL to test
```

---

## 📁 Project Architecture

```
Backend = Layer Cake

┌─────────────────────────────────────┐
│  REST Endpoints (handlers/)         │ ← HTTP layer
├─────────────────────────────────────┤
│  Business Logic (services/)         │ ← Where you add logic
├─────────────────────────────────────┤
│  Database Access (repositories/)    │ ← Query builders
├─────────────────────────────────────┤
│  Data Models (models/)              │ ← GORM models
├─────────────────────────────────────┤
│  PostgreSQL Database                │ ← Data storage
└─────────────────────────────────────┘

Data Flow:
Users → REST API → Handlers → Services → Repositories → Database
                                        → Response format (DTOs)
```

---

## 📋 Essential Files

| File | Purpose | Status |
|------|---------|--------|
| `.env` | Configuration | Create yourself |
| `cmd/server/main.go` | Entry point | ✅ Working |
| `internal/handlers/` | HTTP handlers | ✅ Ready |
| `internal/services/` | Business logic | 80% ✅ |
| `internal/dto/` | Response format | 100% ✅ |
| `internal/routes/routes.go` | URL routing | ✅ Working |
| `Makefile` | Build commands | ✅ Use this |
| `go.mod` | Dependencies | ✅ Locked |

---

## ⚙️ Configuration (.env)

```env
# Server
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=aviator_db

# JWT
JWT_SECRET=change-me-in-production

# Game
GAME_TICK_MS=100
COUNTDOWN_SECONDS=5
```

---

## 🛠️ Make Commands

```bash
make run          # Start server (development)
make build        # Compile binary
make test         # Run all tests
make test-cover   # Tests + coverage report
make fmt          # Format code
make vet          # Check for issues
make clean        # Remove build artifacts
make deps         # Download dependencies
make db-setup     # Create PostgreSQL database
```

---

## 🧪 Testing Endpoints

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@test.com",
    "password": "Pass123!"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@test.com",
    "password": "Pass123!"
  }'
```

### Get Profile (using token)
```bash
# Set TOKEN from login response
$TOKEN = "eyJhbGciOi..."

curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/profile
```

---

## 🐛 Debugging

### Check if Server is Running
```bash
curl http://localhost:8080/api/v1/health
# or
Test-NetConnection localhost -Port 8080
```

### Check Database Connection
```bash
psql -U postgres -d aviator_db -c "SELECT * FROM users LIMIT 1;"
```

### View Server Logs
```bash
# Logs appear in terminal where you ran: make run
# Look for ERROR, WARN, or DEBUG lines
```

### Test Database Manually
```bash
# Connect to database
psql -U postgres -d aviator_db

# List tables
\dt

# Check users
SELECT * FROM users;

# Check wallet balance
SELECT * FROM wallets;

# Exit
\q
```

---

## 📊 Response Format

### Success Response
```json
{
  "status": "success",
  "data": {
    "id": "00000000-0000-5000-8000-000000000001",
    "name": "John",
    "balance": 1000.00,
    "created_at": "2024-01-01T10:00:00Z"
  },
  "timestamp": "2024-01-01T10:30:00Z"
}
```

### Error Response
```json
{
  "status": "error",
  "error": {
    "code": 401,
    "message": "Invalid credentials"
  },
  "timestamp": "2024-01-01T10:30:00Z"
}
```

### Paginated Response
```json
{
  "status": "success",
  "data": {
    "items": [...],
    "pagination": {
      "total": 100,
      "page": 1,
      "limit": 10,
      "pages": 10
    }
  },
  "timestamp": "2024-01-01T10:30:00Z"
}
```

---

## 🔑 Important Concepts

### IDs: Database vs API
```
Database: uint (1, 2, 3...)
API:      UUID strings ("00000000-0000-5000-8000-000000000001")

Convert with:
utils.UintToUUID(1)                    // → UUID string
utils.UUIDStringToUint("00000000...")  // → uint
```

### Timestamps
```
Format: RFC3339 (ISO 8601)
Example: "2024-04-07T10:30:00Z"

Create with:
time.Now().Format(time.RFC3339)
```

### Error Codes
```
200 = OK
201 = Created
400 = Bad Request
401 = Unauthorized
404 = Not Found
409 = Conflict (email exists)
422 = Validation failed
500 = Server error
```

---

## 📚 Documentation Files

Created for you in project root:

| File | Purpose |
|------|---------|
| `API_IMPLEMENTATION_COMPLETE.md` | **START HERE** - Full API spec implementation |
| `BACKEND_MANAGEMENT_GUIDE.md` | How to develop, test, deploy |
| `COMPLETION_CHECKLIST.md` | Remaining work with step-by-step |
| `BACKEND_QUICK_REFERENCE.md` | This file - quick lookup |
| `API_SPECIFICATION.md` | Your original spec (reference) |
| `README.md` | Project overview |
| `ARCHITECTURE.md` | System design |

---

## ✅ What's Complete (70%)

- ✅ All 8 DTOs updated (API request/response format)
- ✅ UUID converter utility created
- ✅ AuthService (register, login, verify)
- ✅ ProfileService (get, update, statistics)
- ✅ WalletService (get wallet)
- ✅ BetService (place, cashout, history)
- ✅ All 14 REST endpoints defined
- ✅ Response wrapper format standardized
- ⚠️ WebSocket events defined (broadcaster needs implementation)

---

## ⚠️ What Needs Work (30%)

**High Priority** (Must do before production):
1. Handler ID Conversion (1-2 hrs) - Parse UUID strings
2. GameService UUID Update (1-2 hrs) - Same format as BetService
3. TransactionService Review (30 mins) - Verify DTOs
4. WebSocket Broadcaster (2-3 hrs) - Implement event routing
5. Integration Tests (3-4 hrs) - End-to-end testing

**Low Priority** (Nice to have):
- Rate limiting
- Redis caching
- Admin endpoints

---

## 🔄 Common Workflows

### Adding a New Endpoint
1. Create DTO in `internal/dto/`
2. Add method to service in `internal/services/`
3. Add handler method in `internal/handlers/`
4. Add route in `internal/routes/routes.go`
5. Test with cURL
6. Commit changes

### Fixing a Bug
1. Identify the issue (check logs)
2. Add debug output/logging
3. Test locally with cURL
4. Verify database state
5. Remove debug code
6. Run tests: `make test`
7. Commit fix

### Deploying Changes
1. `make fmt` - Format
2. `make test` - Test
3. `make build` - Build
4. Copy binary to server
5. Run: `./aviator_backend`
6. Test endpoints
7. Monitor logs

---

## 📞 Problem Solving

| Problem | Solution |
|---------|----------|
| Port 8080 in use | `lsof -i :8080 \| grep LISTEN \| awk '{print $2}' \| xargs kill` |
| DB connection error | Check PostgreSQL running: `psql --version` |
| Tests failing | `make clean && go mod tidy && make test` |
| JWT error | Set `JWT_SECRET` in `.env` |
| UUID conversion panic | Add error handling: `if err != nil { ... }` |
| WebSocket not connecting | Check broadcast running in main.go |

---

## 📞 PowerShell Helpers

```powershell
# Set environment variables
$env:APP_PORT = "8080"
$env:JWT_SECRET = "dev-secret"

# Check what's using a port
Get-NetTcpConnection -LocalPort 8080

# Kill process on port
Get-Process -Id (Get-NetTcpConnection -LocalPort 8080).OwningProcess | Stop-Process

# Build and run
go build -o bin/app.exe cmd/server/main.go
.\bin\app.exe

# Test endpoint
$response = curl http://localhost:8080/api/v1/health
$response | ConvertFrom-Json | Format-Table
```

---

## 🚀 Flutter Integration Points

Your Flutter app should:

1. **REST Endpoints** - Hit `http://localhost:8080/api/v1/{endpoint}`
2. **Authentication** - Send JWT in header: `Authorization: Bearer {token}`
3. **WebSocket** - Connect to `ws://localhost:8080/ws`
4. **Error Handling** - Parse error.code and error.message
5. **Pagination** - Use page, limit, pages from pagination metadata
6. **IDs** - Treat all IDs as strings (UUIDs)

---

## 📈 Performance Tips

- Use indexes on frequently searched columns
- Cache game configurations
- Batch database operations
- Use connection pooling (GORM does this)
- Monitor query performance: `make test-cover`

---

## 🔒 Security Checklist

- [ ] Change JWT_SECRET before production
- [ ] Use HTTPS in production
- [ ] Validate all inputs (already in DTOs)
- [ ] Hash passwords (using bcrypt ✅)
- [ ] Check user ownership of resources
- [ ] Rate limiting enabled
- [ ] CORS configured properly
- [ ] No secrets in git history

---

## 📋 Useful Queries

```bash
# Count users
psql -U postgres -d aviator_db -c "SELECT COUNT(*) as total_users FROM users;"

# View recent transactions
psql -U postgres -d aviator_db -c "SELECT * FROM transactions ORDER BY created_at DESC LIMIT 10;"

# Check wallet balance
psql -U postgres -d aviator_db -c "SELECT users.name, wallets.balance FROM users JOIN wallets ON users.id = wallets.user_id LIMIT 5;"

# Reset user balance
psql -U postgres -d aviator_db -c "UPDATE wallets SET balance = 1000 WHERE user_id = 1;"
```

---

## 🎯 Next Steps

1. **Read** `API_IMPLEMENTATION_COMPLETE.md` (5 mins)
2. **Run** `make run` and test health endpoint (5 mins)
3. **Complete** Task 1: Handler ID Conversion (1-2 hrs)
4. **Complete** Task 2: GameService Update (1-2 hrs)
5. **Complete** Task 3: TransactionService Review (30 mins)
6. **Test** All endpoints manually (1 hr)
7. **Connect** Flutter frontend (depends on your pace)

**Total Time to Completion**: 6-10 hours

---

## 🆘 Need Help?

Check in this order:
1. This file (BACKEND_QUICK_REFERENCE.md)
2. BACKEND_MANAGEMENT_GUIDE.md
3. API_IMPLEMENTATION_COMPLETE.md
4. COMPLETION_CHECKLIST.md
5. Error logs from `make run`
6. Git history for similar implementations

---

## 💾 Saving Your Work

```bash
# Initialize git (if not done)
git init

# Add all changes
git add .

# Commit
git commit -m "feat: Update API implementation to 70% complete"

# View history
git log --oneline
```

---

## 📞 Quick Commands Reference

```bash
# Development
make run              # Start server
make fmt              # Format code
make test             # Run tests

# Database
make db-setup         # Create DB
make db-reset         # Clear DB

# Building
make build            # Compile
make clean            # Remove builds
make deps             # Download deps

# Testing
curl http://localhost:8080/api/v1/health
psql -U postgres -d aviator_db -c "SELECT 1;"
```

---

**Last Updated**: 2024-04-07  
**Created For**: Aviator Betting Game Backend Management  
**Quick Reference Version**: 1.0  
**Status**: Ready to Use  
