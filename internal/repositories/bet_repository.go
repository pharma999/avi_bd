package repositories

import (
	"avi_bd/internal/models"

	"gorm.io/gorm"
)

// BetRepository handles bet database operations
type BetRepository struct {
	db *gorm.DB
}

// NewBetRepository creates a new bet repository
func NewBetRepository(db *gorm.DB) *BetRepository {
	return &BetRepository{db: db}
}

// Create creates a new bet
func (r *BetRepository) Create(bet *models.Bet) error {
	return r.db.Create(bet).Error
}

// GetByID retrieves a bet by ID
func (r *BetRepository) GetByID(id uint) (*models.Bet, error) {
	var bet models.Bet
	err := r.db.First(&bet, id).Error
	if err != nil {
		return nil, err
	}
	return &bet, nil
}

// GetByUserAndRound retrieves a bet by user and round
func (r *BetRepository) GetByUserAndRound(userID, roundID uint) (*models.Bet, error) {
	var bet models.Bet
	err := r.db.Where("user_id = ? AND round_id = ?", userID, roundID).First(&bet).Error
	if err != nil {
		return nil, err
	}
	return &bet, nil
}

// GetByRoundID retrieves all bets for a round
func (r *BetRepository) GetByRoundID(roundID uint) ([]models.Bet, error) {
	var bets []models.Bet
	err := r.db.Where("round_id = ?", roundID).Find(&bets).Error
	return bets, err
}

// GetByUserID retrieves all bets by user ID with pagination
func (r *BetRepository) GetByUserID(userID uint, limit, offset int) ([]models.Bet, int64, error) {
	var bets []models.Bet
	var count int64

	err := r.db.Where("user_id = ?", userID).Model(&models.Bet{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&bets).Error
	if err != nil {
		return nil, 0, err
	}

	return bets, count, nil
}

// Update updates a bet
func (r *BetRepository) Update(bet *models.Bet) error {
	return r.db.Save(bet).Error
}

// GetActiveBetsByRound retrieves all active bets for a round
func (r *BetRepository) GetActiveBetsByRound(roundID uint) ([]models.Bet, error) {
	var bets []models.Bet
	err := r.db.Where("round_id = ? AND result = ?", roundID, "pending").Find(&bets).Error
	return bets, err
}
