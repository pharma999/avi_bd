package models

import (
	"time"

	"gorm.io/gorm"
)

// GameRound represents a single game round
type GameRound struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RoundCode       string         `gorm:"uniqueIndex;size:255" json:"round_code"`
	StartTime       time.Time      `json:"start_time"`
	CrashMultiplier float64        `gorm:"type:decimal(10,2)" json:"crash_multiplier"`
	Status          string         `gorm:"index;size:50" json:"status"` // waiting, running, crashed, ended
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Bets []Bet `gorm:"foreignKey:RoundID" json:"bets,omitempty"`
}

// TableName specifies the table name for GameRound
func (GameRound) TableName() string {
	return "game_rounds"
}
