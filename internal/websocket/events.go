package websocket

import (
	"time"

	"avi_bd/internal/dto"
	"avi_bd/internal/utils"
)

// EventType constants
const (
	EventTypeCountdown        = "countdown"
	EventTypeRoundStart       = "game_start"
	EventTypeMultiplierUpdate = "multiplier_update"
	EventTypeCrash            = "game_crash"
	EventTypeRoundEnd         = "round_end"
	EventTypeBetPlaced        = "bet_placed"
	EventTypeBetCashedOut     = "bet_cashed_out"
	EventTypeError            = "error"
	EventTypePing             = "ping"
	EventTypePong             = "pong"
)

// NewCountdownEvent creates a countdown event
func NewCountdownEvent(roundID uint, secondsRemaining int) *dto.CountdownEvent {
	return &dto.CountdownEvent{
		Type:             EventTypeCountdown,
		RoundID:          utils.UintToUUID(roundID),
		SecondsRemaining: secondsRemaining,
		Timestamp:        time.Now().Format(time.RFC3339),
	}
}

// NewGameStartEvent creates a game start event
func NewGameStartEvent(roundID uint) *dto.GameStartEvent {
	return &dto.GameStartEvent{
		Type:      EventTypeRoundStart,
		RoundID:   utils.UintToUUID(roundID),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// NewMultiplierUpdateEvent creates a multiplier update event
func NewMultiplierUpdateEvent(roundID uint, multiplier float64) *dto.MultiplierUpdateEvent {
	return &dto.MultiplierUpdateEvent{
		Type:             EventTypeMultiplierUpdate,
		RoundID:          utils.UintToUUID(roundID),
		CurrentMultiplier: multiplier,
		Timestamp:        time.Now().Format(time.RFC3339),
	}
}

// NewGameCrashEvent creates a crash event
func NewGameCrashEvent(roundID uint, crashMultiplier float64) *dto.GameCrashEvent {
	return &dto.GameCrashEvent{
		Type:            EventTypeCrash,
		RoundID:         utils.UintToUUID(roundID),
		FinalMultiplier: crashMultiplier,
		Timestamp:       time.Now().Format(time.RFC3339),
	}
}

// NewBetPlacedEvent creates a bet placed event
func NewBetPlacedEvent(roundID uint, amount float64) *dto.BetPlacedEvent {
	return &dto.BetPlacedEvent{
		Type:      EventTypeBetPlaced,
		RoundID:   utils.UintToUUID(roundID),
		Amount:    amount,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// NewBetCashedOutEvent creates a bet cashed out event
func NewBetCashedOutEvent(roundID uint, multiplier float64, winnings float64) *dto.BetCashedOutEvent {
	return &dto.BetCashedOutEvent{
		Type:       EventTypeBetCashedOut,
		RoundID:    utils.UintToUUID(roundID),
		Multiplier: multiplier,
		Winnings:   winnings,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// NewErrorEvent creates an error event
func NewErrorEvent(message string, code int) *dto.ErrorEvent {
	return &dto.ErrorEvent{
		Type:      EventTypeError,
		Message:   message,
		Code:      code,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// NewPingMessage creates a ping message
func NewPingMessage() *dto.PingMessage {
	return &dto.PingMessage{
		Type: EventTypePing,
	}
}

// NewPongMessage creates a pong message
func NewPongMessage() *dto.PongMessage {
	return &dto.PongMessage{
		Type: EventTypePong,
	}
}
