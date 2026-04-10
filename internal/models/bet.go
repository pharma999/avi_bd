package models

import (
	"time"

	"gorm.io/gorm"
)

// Bet represents a user's bet on a game round
type Bet struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"index" json:"user_id"`
	RoundID           uint           `gorm:"index" json:"round_id"`
	Amount            float64        `gorm:"type:decimal(20,2)" json:"amount"`
	AutoCashout       float64        `gorm:"type:decimal(10,2)" json:"auto_cashout"`
	CashoutMultiplier float64        `gorm:"type:decimal(10,2)" json:"cashout_multiplier"`
	Result            string         `gorm:"size:50" json:"result"` // pending, win, loss
	Payout            float64        `gorm:"type:decimal(20,2)" json:"payout"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign keys
	User  *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Round *GameRound `gorm:"foreignKey:RoundID" json:"round,omitempty"`
}

// TableName specifies the table name for Bet
func (Bet) TableName() string {
	return "bets"
}
