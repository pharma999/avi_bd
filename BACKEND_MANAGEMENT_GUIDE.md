# Backend Project Management Guide
## Aviator Betting Game - Go Backend

**Project Location**: `D:\go project\avi_bd`  
**Status**: 70% Complete - Ready for Integration Testing  
**Framework**: Gin + PostgreSQL + WebSocket

---

## 📋 Table of Contents
1. [Project Structure](#project-structure)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Common Tasks](#common-tasks)
5. [Testing & Debugging](#testing--debugging)
6. [Git & Version Control](#git--version-control)
7. [Integration with Flutter](#integration-with-flutter)
8. [Deployment](#deployment)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

---

## 🏗️ Project Structure

```
D:\go project\avi_bd
├── cmd/server/
│   └── main.go                 # Application entry point
├── config/
│   ├── config.go              # Configuration loading
│   └── env.go                 # Environment variables
├── internal/
│   ├── auth/                  # JWT & authentication
│   │   ├── jwt.go
│   │   ├── middleware.go
│   │   └── password.go
│   ├── constants/             # Application constants
│   ├── database/              # Database setup & migrations
│   │   ├── migrations.go
│   │   ├── postgres.go
│   │   └── seed.go
│   ├── dto/                   # Data Transfer Objects ✅ UPDATED
│   │   ├── auth_dto.go
│   │   ├── bet_dto.go
│   │   ├── profile_dto.go
│   │   ├── wallet_dto.go
│   │   ├── transaction_dto.go
│   │   ├── history_dto.go
│   │   ├── response_dto.go
│   │   └── websocket_dto.go
│   ├── game/                  # Game engine (Aviator logic)
│   │   ├── crash_generator.go
│   │   ├── engine.go
│   │   ├── payout_calculator.go
│   │   └── round_manager.go
│   ├── handlers/              # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── bet_handler.go
│   │   ├── profile_handler.go
│   │   ├── wallet_handler.go
│   │   ├── history_handler.go
│   │   └── health_handler.go
│   ├── models/                # Data models (GORM)
│   │   ├── user.go
│   │   ├── wallet.go
│   │   ├── bet.go
│   │   ├── game_round.go
│   │   └── transaction.go
│   ├── repositories/          # Database access layer
│   │   ├── user_repository.go
│   │   ├── wallet_repository.go
│   │   ├── bet_repository.go
│   │   ├── game_round_repository.go
│   │   └── transaction_repository.go
│   ├── services/              # Business logic ✅ PARTIALLY UPDATED
│   │   ├── auth_service.go ✅
│   │   ├── bet_service.go ✅
│   │   ├── game_service.go
│   │   ├── profile_service.go ✅
│   │   ├── transaction_service.go
│   │   └── wallet_service.go ✅
│   ├── routes/                # Route definitions
│   │   └── routes.go
│   ├── utils/                 # Utility functions ✅ NEW
│   │   ├── converter.go ✅ (NEW - UUID conversion)
│   │   ├── logger.go
│   │   ├── response.go
│   │   ├── time.go
│   │   └── validator.go
│   └── websocket/             # WebSocket management
│       ├── broadcaster.go
│       ├── client.go
│       ├── events.go
│       ├── hub.go
├── pkg/
│   └── common/
│       └── errors.go
├── scripts/
│   ├── migrate.sh
│   └── run.sh
├── go.mod                      # Dependencies
├── go.sum
├── Makefile                    # Build commands
├── README.md
├── QUICKSTART.md
├── PROJECT_SUMMARY.md
├── ARCHITECTURE.md
├── FILE_REFERENCE.md
└── API_IMPLEMENTATION_COMPLETE.md ✅ (NEW)
```

**Legend**: ✅ = Recently Updated | ⚠️ = Needs Work

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12+
- Git (optional but recommended)

### Step 1: Setup Database

```bash
cd D:\go project\avi_bd

# Option A: Using psql directly
psql -U postgres -c "CREATE DATABASE aviator_db;"

# Option B: Using the Makefile
make db-setup
```

### Step 2: Configure Environment

Create or update `.env` file in project root:

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
JWT_SECRET=your-super-secret-key-change-in-production

# Game Settings
GAME_TICK_MS=100
COUNTDOWN_SECONDS=5
```

### Step 3: Install Dependencies

```bash
# Download all dependencies
go mod download

# Alternatively
make deps
```

### Step 4: Run the Backend

```bash
# Option A: Direct run (development)
go run cmd/server/main.go

# Option B: Using Makefile
make run

# Option C: Build binary first
make build
./bin/aviator_backend
```

**Expected Output:**
```
✓ Database connected: aviator_db
✓ Migrations completed
✓ Server starting on http://localhost:8080
✓ WebSocket available at ws://localhost:8080/ws
✓ Swagger docs at http://localhost:8080/swagger/index.html
```

---

## 💻 Development Workflow

### Daily Development Cycle

```
1. Start your day
   ├─ Pull latest changes (if team)
   ├─ make deps (if dependencies changed)
   └─ make run (start the backend)

2. Make changes to code
   ├─ Edit files in internal/
   ├─ Save and the app auto-reloads (if using hot-reload)
   └─ Test with curl or Postman

3. Before committing
   ├─ make fmt (format code)
   ├─ make vet (check for errors)
   ├─ make test (run tests)
   └─ make clean (cleanup)

4. Commit and push
```

### Code Organization Rules

**Follow the layer architecture:**

```
Handler (HTTP request)
   ↓
Service (Business logic) ← ADD YOUR LOGIC HERE FIRST
   ↓
Repository (Database access)
   ↓
Model (Data structure)
```

**Example: Adding a new endpoint**
1. Create DTO in `internal/dto/`
2. Add method to service in `internal/services/`
3. Add method to repository (if db access needed)
4. Create handler in `internal/handlers/`
5. Add route in `internal/routes/routes.go`

---

## 📝 Common Tasks

### Task 1: Add a New API Endpoint

**Scenario**: Add GET `/api/v1/me` to get current user info

**Step 1**: Create DTO
```go
// internal/dto/auth_dto.go (add to existing file)
type CurrentUserResponse struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    Email     string  `json:"email"`
    Balance   float64 `json:"balance"`
    CreatedAt string  `json:"created_at"`
}
```

**Step 2**: Add to service
```go
// internal/services/auth_service.go
func (s *AuthService) GetCurrentUser(userID uint) (*dto.CurrentUserResponse, error) {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    wallet, _ := s.walletRepo.GetByUserID(userID)
    
    return &dto.CurrentUserResponse{
        ID:        utils.UintToUUID(user.ID),
        Name:      user.Name,
        Email:     user.Email,
        Balance:   wallet.Balance,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
    }, nil
}
```

**Step 3**: Add handler
```go
// internal/handlers/auth_handler.go
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    userID := auth.GetUserID(c)
    
    user, err := h.authService.GetCurrentUser(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Current user retrieved", user)
}
```

**Step 4**: Add route
```go
// internal/routes/routes.go (add to protected routes)
protected.GET("/me", authHandler.GetCurrentUser)
```

### Task 2: Modify Validation

**Scenario**: Change minimum password length from 8 to 12

```go
// internal/dto/auth_dto.go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=1,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=12,max=255"` // Changed from min=8
}
```

### Task 3: Add Database Field

**Scenario**: Add phone number to User model

**Step 1**: Update model
```go
// internal/models/user.go
type User struct {
    // ... existing fields
    PhoneNumber string `gorm:"column:phone_number;size:20" json:"phone_number"`
}
```

**Step 2**: Create migration
```go
// internal/database/migrations.go
// In AutoMigrate() or create new migration file
type CreateUsersTableMigration struct{}

func (m *CreateUsersTableMigration) Up(db *gorm.DB) error {
    return db.Migrator().AddColumn(&models.User{}, "phone_number")
}
```

**Step 3**: Update DTOs
```go
// internal/dto/profile_dto.go
type ProfileResponse struct {
    // ... existing
    PhoneNumber string `json:"phone_number,omitempty"`
}
```

### Task 4: Fix a Bug in Service

**Example**: Bet not crediting wallet properly

```bash
# 1. Identify the issue
#    - Check logs: make run (look for errors)
#    - Check database: SELECT * FROM wallets WHERE user_id = ?
#    - Add debug logs

# 2. Add debug logging
// internal/services/bet_service.go
func (s *BetService) Cashout(...) {
    // ... existing code
    fmt.Printf("DEBUG: Before credit - Balance: %f\n", currentBalance)
    
    if err := s.walletRepo.IncrementBalance(userID, payout); err != nil {
        return nil, err
    }
    
    newBalance, _ := s.walletRepo.GetBalance(userID)
    fmt.Printf("DEBUG: After credit - Balance: %f\n", newBalance)
}

# 3. Test and verify
make fmt
make test
make run

# 4. Remove debug code once fixed
```

### Task 5: Update Business Logic

**Scenario**: Change auto-cashout min from 1.0 to 1.1

```go
// internal/constants/game_status.go (or wherever constants are)
const (
    MinAutoMultiplier = 1.1 // Changed from 1.0
)

// internal/dto/bet_dto.go
type BetRequest struct {
    // ...
    AutoCashout float64 `json:"auto_cashout" binding:"required,gte=1.1"` // Updated
}

// Update validation in handler if needed
```

---

## 🧪 Testing & Debugging

### Using Makefile Commands

```bash
# Run all tests
make test

# Run tests with coverage report
make test-cover

# Format code
make fmt

# Check for issues
make vet

# Run linter
make lint
```

### Manual Testing with cURL

**Test Authentication Flow:**
```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "TestPass123!"
  }'

# 2. Extract token from response and set it
$TOKEN = "eyJhbGciOiJIUzI1NiI..."

# 3. Get Profile (using token)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/profile

# 4. Get Wallet
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/wallet
```

### Using Postman

1. **Import Collection**:
   - Create new request
   - Set base URL: `http://localhost:8080`
   - Set auth header: `Authorization: Bearer {token}`

2. **Create Environment Variables**:
   - `base_url`: http://localhost:8080
   - `token`: (updated after login)
   - `user_id`: (extracted from responses)

### Debugging

**Enable Debug Mode:**
```go
// In cmd/server/main.go
gin.SetMode(gin.DebugMode) // Already enabled in development
```

**View Detailed Logs:**
```bash
# Run with detailed output
make run 2>&1 | tee debug.log

# Then analyze debug.log
```

**Common Debug Scenarios:**

1. **Database connection fails**
   ```bash
   # Test connection
   psql -h localhost -U postgres -d aviator_db -c "SELECT 1;"
   
   # Check config
   echo Current DB_HOST: $env:DB_HOST
   ```

2. **JWT token invalid**
   ```bash
   # Check JWT_SECRET is set
   echo $env:JWT_SECRET
   
   # Verify token format (should be 3 parts separated by dots)
   ```

3. **Bet not placed**
   ```bash
   # Check database directly
   psql -U postgres -d aviator_db -c "SELECT * FROM bets WHERE user_id = 1;"
   psql -U postgres -d aviator_db -c "SELECT * FROM wallets WHERE user_id = 1;"
   ```

---

## 🔄 Git & Version Control

### Setup Git (if not already done)

```bash
cd D:\go project\avi_bd

# Initialize repo
git init

# Set your identity
git config user.name "Your Name"
git config user.email "your@email.com"

# Create .gitignore
```

### .gitignore Template

Create `.gitignore` in project root:

```
# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Environment
.env
.env.local
.env.*.local

# IDE
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store

# Test coverage
*.coverage
coverage.out

# Temporary
tmp/
temp/
*.log

# Go modules
vendor/
```

### Common Git Commands

```bash
# Check status
git status

# Stage changes
git add .

# Commit
git commit -m "feat: Add UUID conversion for bet responses"

# View history
git log --oneline

# Create branch for feature
git checkout -b feature/websocket-broadcaster

# Merge back to main
git checkout main
git merge feature/websocket-broadcaster
```

### Committing Best Practices

**Good commit message:**
```
feat: Add UUID conversion utility for API responses

- Created internal/utils/converter.go
- Updated all DTOs to use UUID strings
- Updated services to convert uint IDs to UUID format
- Fixes: User IDs now consistent across API
```

**Avoid:**
```
fixed stuff        ← Too vague
updated files      ← Not descriptive
asdf               ← Not helpful
```

---

## 🔗 Integration with Flutter

### Backend Setup for Flutter

1. **Enable CORS (if Flutter runs on different port)**

```go
// In cmd/server/main.go, after creating router
router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})
```

2. **Configure for Development**

Update `.env`:
```env
APP_PORT=8080
DB_HOST=localhost
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://192.168.1.100:3000
```

3. **API Integration from Flutter**

The Flutter app should:
- Use `http` or `dio` package for REST calls
- Connect to WebSocket at `ws://localhost:8080/ws`
- Store JWT token in SharedPreferences
- Refresh token on 401 response (if implemented)

### Common Flutter-to-Backend Workflows

**User Registration:**
```
Flutter App
  ↓ POST /api/v1/register
Backend (AuthService.Register)
  ↓ Create user + wallet
Database
  ↓ Return JWT token
Flutter App (stores token)
```

**Placing Bet:**
```
Flutter App (with token)
  ↓ WebSocket: send place_bet message
Backend (WebSocket Handler)
  ↓ Verify token + validate bet
BetService.PlaceBet
  ↓ Debit wallet + create bet
Database
  ↓ Broadcast bet_placed event
Connected Clients receive update
```

---

## 🚢 Deployment

### Pre-Deployment Checklist

```bash
# 1. Format and clean code
make fmt
make clean

# 2. Run all tests
make test

# 3. Build production binary
make build

# 4. Check binary size
ls -lh bin/aviator_backend

# 5. Create deployment package
mkdir -p deploy
cp bin/aviator_backend deploy/
cp .env.production deploy/.env
```

### Deployment Steps

**Step 1: Prepare Production Database**
```bash
# Create production database
psql -U postgres -c "CREATE DATABASE aviator_prod;"

# Update .env.production
DB_NAME=aviator_prod
JWT_SECRET=<generate-strong-secret>
```

**Step 2: Deploy Binary**
```bash
# Copy files to server
scp -r deploy/ user@server:/opt/aviator/

# SSH into server
ssh user@server

# Navigate to app directory
cd /opt/aviator

# Run migrations
./aviator_backend  # First run creates tables

# Run in background
nohup ./aviator_backend > aviator.log 2>&1 &
```

**Step 3: Setup Reverse Proxy (Nginx)**
```nginx
server {
    listen 80;
    server_name api.aviator.game;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## ✅ Best Practices

### Code Organization

1. **Keep layers separate**
   - Don't call handlers from services
   - Don't access database directly from handlers
   - Keep business logic in services

2. **Use interfaces**
   ```go
   type BetService interface {
       PlaceBet(userID uint, req *dto.BetRequest) (*dto.BetResponse, error)
       Cashout(userID uint, betID uint, multiplier float64) (*dto.CashoutResponse, error)
   }
   ```

3. **Error handling**
   ```go
   // Always check and handle errors
   if err != nil {
       utils.ErrorResponse(c, http.StatusInternalServerError, "Failed", err.Error())
       return
   }
   ```

### Performance

1. **Use database indexes**
   ```go
   type User struct {
       ID    uint   `gorm:"primaryKey"`
       Email string `gorm:"uniqueIndex;index"`
   }
   ```

2. **Cache frequently accessed data**
   - Consider Redis for user sessions
   - Cache game configurations

3. **Batch database operations**
   ```go
   // Bad: N+1 queries
   for _, bet := range bets {
       round := getRound(bet.RoundID) // N queries
   }
   
   // Good: Preload
   db.Preload("Round").Find(&bets)
   ```

### Security

1. **Validate all inputs**
   - Use struct tags: `binding:"required,email"`
   - Validate in services too

2. **Hash sensitive data**
   - Passwords: bcrypt ✅
   - API keys: consider hashing

3. **Protect endpoints**
   - All authenticated endpoints need JWT middleware
   - Check user ownership of resources

4. **Environment variables**
   - Never commit `.env` with real secrets
   - Use `.env.example` template

---

## 🆘 Troubleshooting

### Issue: "Database connection refused"

**Solution:**
```bash
# 1. Check if PostgreSQL is running
psql --version
psql -U postgres -c "SELECT 1;"

# 2. Verify credentials in .env
cat .env | grep DB_

# 3. Check database exists
psql -U postgres -c "\l"

# 4. Create if missing
psql -U postgres -c "CREATE DATABASE aviator_db;"
```

### Issue: "Port 8080 already in use"

**Solution:**
```bash
# PowerShell: Find and kill process
Get-Process -Id (Get-NetTcpConnection -LocalPort 8080).OwningProcess | Stop-Process

# Or use different port
$env:APP_PORT=8081
make run
```

### Issue: "JWT secret not set"

**Solution:**
```bash
# Generate secure secret
# PowerShell:
[Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((Get-Random -Count 32 -InputObject @(0..255) | ForEach-Object {[char]$_}) -join ''))

# Or set to something
$env:JWT_SECRET="your-super-secret-key"
```

### Issue: "Tests failing locally but should pass"

**Solution:**
```bash
# Clean and rebuild
make clean
make deps
make build

# Run specific test
go test -v -run TestBetService ./internal/services/...

# Check test database
psql -U postgres -c "\l" | grep test
```

### Issue: "WebSocket connection refused"

**Solution:**
```bash
# 1. Check server is running
curl http://localhost:8080/health

# 2. Check WebSocket endpoint
# Use websocat tool or browser DevTools Console
wscat -c ws://localhost:8080/ws

# 3. Verify in handlers
# Confirm ./internal/routes/routes.go has:
// router.GET("/ws", handleWebSocket(hub))
```

---

## 📊 Remaining Work (30%)

### High Priority (Complete these first)

- [ ] **Handler ID Conversion** (1-2 hours)
  - Add UUID-to-uint conversion in handlers
  - File: `internal/handlers/bet_handler.go`, `history_handler.go`
  
- [ ] **GameService Update** (1-2 hours)
  - Apply same UUID pattern as BetService
  - File: `internal/services/game_service.go`

- [ ] **TransactionService Review** (30 mins)
  - File: `internal/services/transaction_service.go`

### Medium Priority

- [ ] **WebSocket Broadcaster** (2-3 hours)
  - Implement event routing
  - File: `internal/websocket/broadcaster.go`

- [ ] **Integration Tests** (3-4 hours)
  - Full workflow testing
  - File: Create `internal/*/test.go` files

### Lower Priority

- [ ] Rate limiting implementation
- [ ] Cache layer (Redis)
- [ ] Admin endpoints
- [ ] Advanced analytics

---

## 📖 Quick Reference

```bash
# Start development
make run

# Code quality
make fmt && make vet && make test

# Build for production
make build

# View project structure
ls -R internal/

# Check API docs
curl http://localhost:8080/swagger/index.html

# View logs
tail -f debug.log

# Database access
psql -U postgres -d aviator_db

# Kill running process
lsof -i :8080 | grep LISTEN | awk '{print $2}' | xargs kill -9
```

---

## 🎯 Next Steps

1. **Review the Implementation** - Read `API_IMPLEMENTATION_COMPLETE.md`
2. **Complete Remaining 30%** - Start with high priority items above
3. **Test Thoroughly** - Manual testing + unit tests
4. **Connect Flutter** - Point frontend to http://localhost:8080
5. **Deploy** - Follow deployment section when ready

---

**Created**: 2024-04-07  
**Last Updated**: 2024-04-07  
**Status**: Ready for Development  
**Support**: Check `API_IMPLEMENTATION_COMPLETE.md` for API details
