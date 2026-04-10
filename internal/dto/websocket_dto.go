package dto

// WebSocketEvent represents a websocket event
type WebSocketEvent struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload,omitempty"`
	RoundID   string      `json:"round_id,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Message   string      `json:"message,omitempty"`
	Code      int         `json:"code,omitempty"`
}

// AuthMessage represents WebSocket authentication message
type AuthMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// PlaceBetMessage represents WebSocket bet placement message
type PlaceBetMessage struct {
	Type         string  `json:"type"`
	RoundID      string  `json:"round_id"`
	Amount       float64 `json:"amount"`
	AutoCashout  float64 `json:"auto_cashout"`
}

// CashoutMessage represents WebSocket cashout message
type CashoutMessage struct {
	Type              string  `json:"type"`
	BetID             string  `json:"bet_id"`
	CurrentMultiplier float64 `json:"current_multiplier"`
}

// CountdownEvent represents countdown event sent before each round
type CountdownEvent struct {
	Type             string `json:"type"`
	RoundID          string `json:"round_id"`
	SecondsRemaining int    `json:"seconds_remaining"`
	Timestamp        string `json:"timestamp"`
}

// GameStartEvent represents game_start event per specification
type GameStartEvent struct {
	Type      string `json:"type"`
	RoundID   string `json:"round_id"`
	Timestamp string `json:"timestamp"`
}

// MultiplierUpdateEvent represents multiplier_update event per specification
type MultiplierUpdateEvent struct {
	Type             string  `json:"type"`
	RoundID          string  `json:"round_id"`
	CurrentMultiplier float64 `json:"current_multiplier"`
	Timestamp        string  `json:"timestamp"`
}

// GameCrashEvent represents game_crash event per specification
type GameCrashEvent struct {
	Type            string  `json:"type"`
	RoundID         string  `json:"round_id"`
	FinalMultiplier float64 `json:"final_multiplier"`
	Timestamp       string  `json:"timestamp"`
}

// BetPlacedEvent represents bet_placed event per specification
type BetPlacedEvent struct {
	Type      string  `json:"type"`
	RoundID   string  `json:"round_id"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
}

// BetCashedOutEvent represents bet_cashed_out event per specification
type BetCashedOutEvent struct {
	Type       string  `json:"type"`
	RoundID    string  `json:"round_id"`
	Multiplier float64 `json:"multiplier"`
	Winnings   float64 `json:"winnings"`
	Timestamp  string  `json:"timestamp"`
}

// ErrorEvent represents error event per specification
type ErrorEvent struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp string `json:"timestamp"`
}

// PingMessage represents ping message per specification
type PingMessage struct {
	Type string `json:"type"`
}

// PongMessage represents pong message per specification
type PongMessage struct {
	Type string `json:"type"`
}
