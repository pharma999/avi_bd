package game

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"avi_bd/config"
	"avi_bd/internal/constants"
	"avi_bd/internal/models"
	"avi_bd/internal/repositories"
	"avi_bd/internal/websocket"

	"gorm.io/gorm"
)

// Engine manages the Aviator game loop.
//
// Crash modes:
//   auto   (default) – plane crashes at a randomly generated multiplier.
//   manual           – plane flies indefinitely until an admin triggers ForceCrash().
type Engine struct {
	db           *gorm.DB
	cfg          *config.Config
	hub          *websocket.Hub
	roundManager *RoundManager
	crashGen     *CrashGenerator
	payoutCalc   *PayoutCalculator
	mu           sync.RWMutex
	isRunning    bool

	// manual crash mode
	manualMode   atomic.Bool      // true = manual, false = auto
	crashSignal  chan struct{}     // admin sends here to force a crash
	crashMu      sync.Mutex       // guards crashSignal replacement between rounds
}

// NewEngine creates a new game engine
func NewEngine(db *gorm.DB, cfg *config.Config, hub *websocket.Hub) *Engine {
	e := &Engine{
		db:           db,
		cfg:          cfg,
		hub:          hub,
		roundManager: NewRoundManager(db),
		crashGen:     NewCrashGenerator(cfg),
		payoutCalc:   NewPayoutCalculator(),
		crashSignal:  make(chan struct{}, 1),
	}
	return e
}

// Start starts the game engine loop
func (e *Engine) Start() {
	go e.gameLoop()
}

// SetManualMode switches between auto (false) and manual (true) crash modes.
// Broadcasts a mode-change event so dashboards can update in real time.
func (e *Engine) SetManualMode(manual bool) {
	e.manualMode.Store(manual)
	modeStr := "auto"
	if manual {
		modeStr = "manual"
	}
	e.hub.Broadcast(map[string]interface{}{
		"type": "crash_mode_changed",
		"mode": modeStr,
	})
}

// IsManualMode returns the current crash mode.
func (e *Engine) IsManualMode() bool {
	return e.manualMode.Load()
}

// ForceCrash triggers an immediate crash in manual mode.
// Returns false if not in manual mode or no round is running.
func (e *Engine) ForceCrash() bool {
	if !e.manualMode.Load() {
		return false
	}
	e.crashMu.Lock()
	ch := e.crashSignal
	e.crashMu.Unlock()

	// Non-blocking send — if already signalled, do nothing
	select {
	case ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// gameLoop is the main game loop
func (e *Engine) gameLoop() {
	e.mu.Lock()
	e.isRunning = true
	e.mu.Unlock()

	for e.isRunning {
		// Fresh crash signal channel each round
		e.crashMu.Lock()
		e.crashSignal = make(chan struct{}, 1)
		e.crashMu.Unlock()

		round := e.roundManager.CreateRound(e.db)
		if round == nil {
			time.Sleep(time.Duration(e.cfg.GameTickMs) * time.Millisecond)
			continue
		}

		e.countdownPhase(round)
		e.runningPhase(round)
		e.db.First(round)
		e.roundEndedPhase(round)

		time.Sleep(2 * time.Second)
	}
}

// countdownPhase broadcasts a countdown event every second.
func (e *Engine) countdownPhase(round *models.GameRound) {
	seconds := e.cfg.CountdownSeconds
	for i := seconds; i > 0; i-- {
		e.hub.Broadcast(websocket.NewCountdownEvent(round.ID, i))
		time.Sleep(time.Second)
	}
}

// runningPhase drives the multiplier and handles crash logic.
// In auto mode  → crashes when multiplier reaches the pre-generated crash point.
// In manual mode → flies indefinitely; crashes only when ForceCrash() is called.
func (e *Engine) runningPhase(round *models.GameRound) {
	round.Status = constants.GameStatusRunning
	e.db.Save(round)

	e.hub.Broadcast(websocket.NewGameStartEvent(round.ID))

	// Pre-generate crash multiplier (used in auto mode; ignored in manual)
	autoCrashAt := e.crashGen.GenerateCrash()
	round.CrashMultiplier = autoCrashAt
	e.db.Save(round)

	startTime := time.Now()
	tickDuration := time.Duration(e.cfg.GameTickMs) * time.Millisecond
	growthRate := 0.05 // 5% per second

	// Snapshot the signal channel for this round
	e.crashMu.Lock()
	signalCh := e.crashSignal
	e.crashMu.Unlock()

	var finalMultiplier float64

	for {
		elapsed := time.Since(startTime).Seconds()
		multiplier := math.Round((1.0+elapsed*growthRate)*100) / 100

		e.hub.Broadcast(websocket.NewMultiplierUpdateEvent(round.ID, multiplier))
		e.processAutoCashouts(round, multiplier)

		if e.manualMode.Load() {
			// Manual mode: only crash on admin signal
			select {
			case <-signalCh:
				finalMultiplier = multiplier
				goto crashed
			default:
			}
		} else {
			// Auto mode: crash when multiplier reaches pre-generated point
			if multiplier >= autoCrashAt {
				finalMultiplier = autoCrashAt
				goto crashed
			}
		}

		time.Sleep(tickDuration)
	}

crashed:
	round.CrashMultiplier = finalMultiplier
	round.Status = constants.GameStatusCrashed
	e.db.Save(round)

	e.hub.Broadcast(websocket.NewGameCrashEvent(round.ID, finalMultiplier))
	e.markBetsAsLosses(round)
}

// roundEndedPhase marks the round as ended
func (e *Engine) roundEndedPhase(round *models.GameRound) {
	round.Status = constants.GameStatusEnded
	e.db.Save(round)
}

// processAutoCashouts processes auto cashouts when multiplier reaches auto_cashout
func (e *Engine) processAutoCashouts(round *models.GameRound, currentMultiplier float64) {
	betRepo := repositories.NewBetRepository(e.db)
	walletRepo := repositories.NewWalletRepository(e.db)

	activeBets, _ := betRepo.GetActiveBetsByRound(round.ID)

	for _, bet := range activeBets {
		if bet.AutoCashout > 0 && currentMultiplier >= bet.AutoCashout && bet.Result == constants.BetResultPending {
			payout := e.payoutCalc.Calculate(bet.Amount, currentMultiplier)

			bet.Result = constants.BetResultWin
			bet.Payout = payout
			bet.CashoutMultiplier = currentMultiplier
			betRepo.Update(&bet)
			walletRepo.IncrementBalance(bet.UserID, payout)

			e.hub.Broadcast(websocket.NewBetCashedOutEvent(round.ID, currentMultiplier, payout))
		}
	}
}

// markBetsAsLosses marks all pending bets in a round as losses
func (e *Engine) markBetsAsLosses(round *models.GameRound) {
	betRepo := repositories.NewBetRepository(e.db)
	activeBets, _ := betRepo.GetActiveBetsByRound(round.ID)
	for _, bet := range activeBets {
		if bet.Result == constants.BetResultPending {
			bet.Result = constants.BetResultLoss
			bet.Payout = 0
			betRepo.Update(&bet)
		}
	}
}

// Stop stops the game engine
func (e *Engine) Stop() {
	e.mu.Lock()
	e.isRunning = false
	e.mu.Unlock()
}

// IsRunning returns whether the engine is running
func (e *Engine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isRunning
}
