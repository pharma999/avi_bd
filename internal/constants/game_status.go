package constants

// Game status constants
const (
	GameStatusWaiting = "waiting"
	GameStatusRunning = "running"
	GameStatusCrashed = "crashed"
	GameStatusEnded   = "ended"
)

// Bet status constants per specification
const (
	BetStatusActive    = "active"
	BetStatusWon       = "won"
	BetStatusLost      = "lost"
	BetStatusCashedOut = "cashed_out"
	BetStatusFailed    = "failed"
)

// Legacy bet result constants (for backward compatibility during refactoring)
const (
	BetResultPending  = "pending"  // Maps to BetStatusActive
	BetResultWin      = "won"      // Maps to BetStatusWon
	BetResultLoss     = "loss"     // Maps to BetStatusLost
	BetResultCashedOut = "cashed_out" // New
)
