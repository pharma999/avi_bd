package dto

// GameRoundItem represents a single game round in history per specification
type GameRoundItem struct {
	ID              string  `json:"id"`
	FinalMultiplier float64 `json:"final_multiplier"`
	CrashedAt       float64 `json:"crashed_at"`
	StartedAt       string  `json:"started_at"`
	EndedAt         string  `json:"ended_at"`
	PlayersCount    int     `json:"players_count"`
}

// GameHistoryResponse represents game history response per specification
type GameHistoryResponse struct {
	Rounds []GameRoundItem `json:"rounds"`
}

// GameRoundResponse represents a single game round per specification
type GameRoundResponse struct {
	ID              string   `json:"id"`
	Status          string   `json:"status"` // waiting, running, crashed
	CurrentMultiplier float64  `json:"current_multiplier"`
	FinalMultiplier *float64 `json:"final_multiplier,omitempty"`
	CrashedAt       *float64 `json:"crashed_at,omitempty"`
	CountdownSeconds int      `json:"countdown_seconds"`
	StartedAt       *string  `json:"started_at,omitempty"`
	EndedAt         *string  `json:"ended_at,omitempty"`
	CreatedAt       string   `json:"created_at"`
}

// CurrentGameResponse represents the current game state
type CurrentGameResponse struct {
	ID                string  `json:"id"`
	Status            string  `json:"status"`
	CurrentMultiplier float64 `json:"current_multiplier"`
	CountdownSeconds  int     `json:"countdown_seconds"`
	StartedAt         *string `json:"started_at,omitempty"`
}
