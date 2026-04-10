package repositories

import (
	"avi_bd/internal/models"

	"gorm.io/gorm"
)

// GameRoundRepository handles game round database operations
type GameRoundRepository struct {
	db *gorm.DB
}

// NewGameRoundRepository creates a new game round repository
func NewGameRoundRepository(db *gorm.DB) *GameRoundRepository {
	return &GameRoundRepository{db: db}
}

// Create creates a new game round
func (r *GameRoundRepository) Create(round *models.GameRound) error {
	return r.db.Create(round).Error
}

// GetByID retrieves a game round by ID
func (r *GameRoundRepository) GetByID(id uint) (*models.GameRound, error) {
	var round models.GameRound
	err := r.db.First(&round, id).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

// GetByRoundCode retrieves a game round by code
func (r *GameRoundRepository) GetByRoundCode(code string) (*models.GameRound, error) {
	var round models.GameRound
	err := r.db.Where("round_code = ?", code).First(&round).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

// GetLatestRound retrieves the latest game round
func (r *GameRoundRepository) GetLatestRound() (*models.GameRound, error) {
	var round models.GameRound
	err := r.db.Order("created_at DESC").First(&round).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

// Update updates a game round
func (r *GameRoundRepository) Update(round *models.GameRound) error {
	return r.db.Save(round).Error
}

// GetByStatus retrieves game rounds by status
func (r *GameRoundRepository) GetByStatus(status string) ([]models.GameRound, error) {
	var rounds []models.GameRound
	err := r.db.Where("status = ?", status).Find(&rounds).Error
	return rounds, err
}

// GetAll retrieves all game rounds with pagination
func (r *GameRoundRepository) GetAll(limit, offset int) ([]models.GameRound, int64, error) {
	var rounds []models.GameRound
	var count int64

	err := r.db.Model(&models.GameRound{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rounds).Error
	if err != nil {
		return nil, 0, err
	}

	return rounds, count, nil
}
