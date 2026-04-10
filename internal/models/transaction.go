package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Type      string         `gorm:"index;size:50" json:"type"` // credit, debit, win, loss, deposit, withdrawal, bet, cashout
	Amount    float64        `gorm:"type:decimal(20,2)" json:"amount"`
	Reference string         `gorm:"size:255" json:"reference"` // bet_id, round_id, etc.
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}
