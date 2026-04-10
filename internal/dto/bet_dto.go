package dto

// BetRequest represents a bet placement request per specification
type BetRequest struct {
	RoundID     string  `json:"round_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	AutoCashout float64 `json:"auto_cashout" binding:"required,gte=1"`
}

// BetResponse represents a bet in response per specification
type BetResponse struct {
	ID             string   `json:"id"`
	UserID         string   `json:"user_id"`
	RoundID        string   `json:"round_id"`
	Amount         float64  `json:"amount"`
	AutoCashout    float64  `json:"auto_cashout"`
	Status         string   `json:"status"` // active, won, lost, cashed_out, failed
	PlacedAt       string   `json:"placed_at"`
	EndedAt        *string  `json:"ended_at,omitempty"`
	CashoutMultiplier *float64 `json:"cashout_multiplier,omitempty"`
	FinalMultiplier *float64 `json:"final_multiplier,omitempty"`
	Winnings       *float64 `json:"winnings,omitempty"`
}

// CashoutRequest represents a cashout request per specification
type CashoutRequest struct {
	BetID             string  `json:"bet_id" binding:"required"`
	CurrentMultiplier float64 `json:"current_multiplier" binding:"required,gte=1"`
}

// CashoutResponse represents a cashout response per specification
type CashoutResponse struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	Status            string  `json:"status"`
	CashoutMultiplier float64 `json:"cashout_multiplier"`
	Winnings          float64 `json:"winnings"`
	CashedOutAt       string  `json:"cashed_out_at"`
}

// BetHistoryItem represents a single bet in history per specification
type BetHistoryItem struct {
	ID                string   `json:"id"`
	RoundID           string   `json:"round_id"`
	Amount            float64  `json:"amount"`
	Status            string   `json:"status"`
	FinalMultiplier   *float64 `json:"final_multiplier,omitempty"`
	CashoutMultiplier *float64 `json:"cashout_multiplier,omitempty"`
	Winnings          *float64 `json:"winnings,omitempty"`
	PlacedAt          string   `json:"placed_at"`
	EndedAt           *string  `json:"ended_at,omitempty"`
}

// BetHistoryResponse represents bet history response per specification
type BetHistoryResponse struct {
	Bets []BetHistoryItem `json:"bets"`
}
