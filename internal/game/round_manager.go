package game

import (
	"fmt"
	"time"

	"avi_bd/internal/constants"
	"avi_bd/internal/models"

	"gorm.io/gorm"
)

// RoundManager manages game rounds
type RoundManager struct {
	db *gorm.DB
}

// NewRoundManager creates a new round manager
func NewRoundManager(db *gorm.DB) *RoundManager {
	return &RoundManager{db: db}
}

// CreateRound creates a new game round
func (rm *RoundManager) CreateRound(db *gorm.DB) *models.GameRound {
	roundCode := rm.generateRoundCode()

	round := &models.GameRound{
		RoundCode: roundCode,
		StartTime: time.Now(),
		Status:    constants.GameStatusWaiting,
	}

	if err := db.Create(round).Error; err != nil {
		return nil
	}

	return round
}

// generateRoundCode generates a unique round code
func (rm *RoundManager) generateRoundCode() string {
	return fmt.Sprintf("ROUND_%d", time.Now().UnixNano())
}

// GetRound retrieves a round by ID
func (rm *RoundManager) GetRound(roundID uint) (*models.GameRound, error) {
	var round models.GameRound
	err := rm.db.First(&round, roundID).Error
	return &round, err
}
