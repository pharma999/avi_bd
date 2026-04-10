package models

import (
	"time"

	"gorm.io/gorm"
)

// Wallet represents a user's wallet/balance
type Wallet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex" json:"user_id"`
	Balance   float64        `gorm:"type:decimal(20,2)" json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Wallet
func (Wallet) TableName() string {
	return "wallets"
}
