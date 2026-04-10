# Completion Checklist - Aviator Backend

## 📋 Overview
**Current Status**: 70% Complete  
**Remaining**: 30% (High Priority Items)  
**Estimated Time**: 6-10 hours  
**Priority**: Complete before Flutter integration

---

## 🔴 HIGH PRIORITY (Must Complete First)

### Task 1: Handler ID Conversion
**Status**: ❌ NOT STARTED  
**Files**: `internal/handlers/bet_handler.go`, `internal/handlers/history_handler.go`  
**Time**: 1-2 hours  
**Difficulty**: Easy  

**What to do**:
- [ ] In bet_handler.go, add UUID string parsing to betID parameters
- [ ] Update GetBetHistory to handle UUID strings
- [ ] In history_handler.go, update round_id parsing
- [ ] Test all endpoints work with UUID format

**Example Changes**:
```go
// Before
betID := c.Param("bet_id") // This is a string "00000000-0000-5000-8000-000000000010"

// After
betIDStr := c.Param("bet_id")
betID, err := utils.UUIDStringToUint(betIDStr)
if err != nil {
    utils.ErrorResponse(c, http.StatusBadRequest, "Invalid bet ID", nil)
    return
}
```

**Verification**:
```bash
# After completing:
make fmt
make test
make run
# Test placing a bet and getting history
```

---

### Task 2: GameService UUID Update
**Status**: ❌ NOT STARTED  
**File**: `internal/services/game_service.go`  
**Time**: 1-2 hours  
**Difficulty**: Easy-Medium  

**What to do**:
- [ ] Add UUID conversion to all response DTOs
- [ ] Update timestamps to RFC3339 format
- [ ] Follow same pattern as BetService
- [ ] Import utils package
- [ ] Add time import

**Pattern to Follow**:
```go
// From BetService (GOOD EXAMPLE)
return &dto.BetResponse{
    ID:       utils.UintToUUID(bet.ID),
    UserID:   utils.UintToUUID(bet.UserID),
    RoundID:  utils.UintToUUID(bet.RoundID),
    PlacedAt: bet.CreatedAt.Format(time.RFC3339),
}
```

**Verification**:
```bash
make test
# Test game history endpoint should return properly formatted responses
```

---

### Task 3: TransactionService Review
**Status**: ⚠️ NEEDS REVIEW  
**File**: `internal/services/transaction_service.go`  
**Time**: 30 mins - 1 hour  
**Difficulty**: Easy  

**What to do**:
- [ ] Review all transaction responses
- [ ] Ensure UUIDs are used for all IDs
- [ ] Verify RFC3339 timestamps
- [ ] Check that DTOs match responses
- [ ] Test GetTransactionHistory endpoint

**Verification**:
```bash
# After review:
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/transactions/history
# Should return properly formatted JSON with UUID IDs and RFC3339 timestamps
```

---

## 🟡 MEDIUM PRIORITY (Important for Full Functionality)

### Task 4: WebSocket Broadcaster Implementation
**Status**: ⚠️ PARTIALLY COMPLETE (Events defined, need broadcaster)  
**Files**: `internal/websocket/broadcaster.go`, `internal/websocket/hub.go`  
**Time**: 2-3 hours  
**Difficulty**: Medium  

**What to do**:
- [ ] Implement game event broadcasting
- [ ] Connect to game engine for updates
- [ ] Send multiplier updates every ~100ms
- [ ] Broadcast bet_placed events
- [ ] Broadcast game_crash events
- [ ] Handle client disconnections
- [ ] Test WebSocket connections

**Events to Implement**:
- game_start
- multiplier_update
- game_crash
- bet_placed
- bet_cashed_out
- error

**Testing**:
```bash
# Use WebSocket client
wscat -c ws://localhost:8080/ws

# Should receive: game_start, multiplier_update, game_crash events
```

---

### Task 5: Integration Tests
**Status**: ❌ NOT STARTED  
**Create Files**: `internal/handlers/*_test.go`, `internal/services/*_test.go`  
**Time**: 3-4 hours  
**Difficulty**: Medium-Hard  

**What to test**:
- [ ] User registration flow
- [ ] Login and token generation
- [ ] Profile retrieval and update
- [ ] Wallet operations
- [ ] Bet placement
- [ ] Bet cashout
- [ ] History retrieval with pagination
- [ ] Error cases (unauthorized, not found, etc.)

**Example Test**:
```go
// internal/handlers/auth_handler_test.go
func TestRegister(t *testing.T) {
    // Setup
    // Call endpoint
    // Assert response
}
```

**Run Tests**:
```bash
make test
make test-cover  # With coverage report
```

---

## 🟢 LOWER PRIORITY (Nice to Have)

### Task 6: Rate Limiting
**Status**: ⚠️ NOT IMPLEMENTED  
**Estimated Time**: 2 hours  
**Difficulty**: Medium  

**Per Specification**:
- 1000 requests/hour per IP
- 500 requests/hour per user
- Return `X-RateLimit-Remaining` header

---

### Task 7: Cache Layer (Redis)
**Status**: ❌ NOT STARTED  
**Estimated Time**: 3-4 hours  
**Difficulty**: Hard  

**What to cache**:
- User sessions
- Game statistics
- Frequently accessed profiles

---

### Task 8: Admin Endpoints
**Status**: ❌ NOT STARTED  
**Estimated Time**: 2-3 hours  
**Difficulty**: Medium  

**Endpoints to add**:
- GET /api/v1/admin/users (paginated)
- GET /api/v1/admin/stats
- POST /api/v1/admin/users/{id}/reset-balance
- DELETE /api/v1/admin/users/{id}

---

## ✅ COMPLETION CHECKLIST

### Pre-Completion
- [ ] Read through BACKEND_MANAGEMENT_GUIDE.md
- [ ] Read through API_IMPLEMENTATION_COMPLETE.md
- [ ] Understand the project structure
- [ ] Run `make run` and test health endpoint

### High Priority Completion
- [ ] Task 1: Handler ID Conversion (1-2 hrs)
- [ ] Task 2: GameService UUID Update (1-2 hrs)
- [ ] Task 3: TransactionService Review (30 mins)
- [ ] All high priority tasks passing `make test`
- [ ] Manual testing of all endpoints

### Medium Priority Completion
- [ ] Task 4: WebSocket Broadcaster (2-3 hrs)
- [ ] Task 5: Integration Tests (3-4 hrs)
- [ ] All tests passing
- [ ] Full workflow tested end-to-end

### Final Verification
- [ ] `make fmt` - Code is formatted
- [ ] `make vet` - No issues found
- [ ] `make test` - All tests pass
- [ ] `make build` - Binary builds successfully
- [ ] Manually test in Postman or cURL
- [ ] Test with Flutter app

### Documentation
- [ ] Update README.md with latest info
- [ ] Comment any complex functions
- [ ] Document any new environment variables
- [ ] Create API usage guide for Flutter team

---

## 🎯 Task Dependencies

```
Legend: ← depends on

Task 1 (Handler ID Conversion)       ← Nothing (START HERE)
  ↓
Task 2 (GameService UUID Update)     ← Task 1
  ↓
Task 3 (TransactionService Review)   ← Task 1, 2
  ↓
Task 4 (WebSocket Broadcaster)       ← Task 1, 2, 3 (All HIGH PRIORITY)
  ↓
Task 5 (Integration Tests)           ← Task 1, 2, 3, 4

Tasks 6-8 (Low Priority)             ← Can be done in parallel with Tasks 1-3
```

---

## 📊 Progress Tracking

### Week 1 Goals
- [x] API specification document reviewed
- [x] All DTOs updated
- [x] Services partially updated (Auth, Profile, Wallet, Bet)
- [x] UUID converter utility created
- [ ] All remaining service updates (Task 2, 3)
- [ ] Handler ID conversions (Task 1)

### Week 2 Goals
- [ ] All 3 high priority tasks complete
- [ ] Manual testing complete
- [ ] WebSocket broadcaster implementation started (Task 4)

### Week 3 Goals
- [ ] WebSocket broadcaster complete
- [ ] Integration tests written and passing (Task 5)
- [ ] Ready for Flutter integration

---

## 🔄 Workflow for Each Task

```
1. Read the task description
2. Open the relevant file
3. Identify what needs to change
4. Make the changes
5. Run: make fmt
6. Run: make test
7. Verify manually
8. Mark as complete
9. Move to next task
```

---

## 📞 Common Issues & Solutions

### Issue: "Test fails with database error"
```bash
# Solution: Reset database
make db-reset
make run
make test
```

### Issue: "UUID conversion panic"
```bash
# Add error handling:
betID, err := utils.UUIDStringToUint(betIDStr)
if err != nil {
    utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
    return
}
```

### Issue: "WebSocket connection timeout"
```bash
# Check broadcaster is running:
# In hub.go, verify Run() method is called in main.go
go hub.Run()  // Should exist in cmd/server/main.go
```

---

## 📝 Commit Messages Template

After completing each task:

```bash
# Task 1
git commit -m "feat: Add UUID conversion to handlers

- Updated bet_handler.go to parse UUID strings
- Updated history_handler.go for UUID format
- All handler endpoints now work with UUID IDs
- Tests passing"

# Task 2
git commit -m "feat: Update GameService with UUID format

- Converted all IDs to UUID strings
- Updated timestamps to RFC3339 format
- Applied same pattern as BetService
- All game endpoints tested"

# Task 3
git commit -m "fix: Review and complete TransactionService

- Verified all transaction responses use UUIDs
- Ensured RFC3339 timestamp formatting
- Transaction history endpoint validated
- Tests passing"
```

---

## 🎓 Learning Resources

**Within Project**:
- `API_IMPLEMENTATION_COMPLETE.md` - API reference
- `BACKEND_MANAGEMENT_GUIDE.md` - Development guide
- `ARCHITECTURE.md` - System architecture
- `FILE_REFERENCE.md` - File descriptions
- Example files: BetService (✅ completed), AuthService (✅ completed)

**External**:
- Go documentation: https://golang.org/doc/
- Gin documentation: https://github.com/gin-gonic/gin
- GORM documentation: https://gorm.io/
- JSON Web Token: https://jwt.io/

---

## ✨ Final Checklist

When you've completed all high priority tasks:

- [ ] Database connection working
- [ ] All 14 REST endpoints return correct format
- [ ] UUIDs used consistently across API
- [ ] RFC3339 timestamps on all responses
- [ ] Pagination working on list endpoints
- [ ] JWT authentication working
- [ ] Error responses properly formatted
- [ ] WebSocket events defined and partially tested
- [ ] Documentation updated
- [ ] Ready for Flutter integration
- [ ] Ready for deployment

**Estimated Time to Complete All**: 8-12 hours  
**Recommended Pace**: 2-3 hours/day

---

**Last Updated**: 2024-04-07  
**Status**: Ready to Start Implementation  
**Next Step**: Start with Task 1: Handler ID Conversion
