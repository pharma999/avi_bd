package repositories

import (
	"avi_bd/internal/models"

	"gorm.io/gorm"
)

// TransactionRepository handles transaction database operations
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetByUserID retrieves transactions by user ID with pagination
func (r *TransactionRepository) GetByUserID(userID uint, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var count int64

	err := r.db.Where("user_id = ?", userID).Model(&models.Transaction{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

// GetByType retrieves transactions by type
func (r *TransactionRepository) GetByType(userID uint, txType string, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var count int64

	err := r.db.Where("user_id = ? AND type = ?", userID, txType).Model(&models.Transaction{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ? AND type = ?", userID, txType).Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

// GetByReference retrieves transactions by reference
func (r *TransactionRepository) GetByReference(reference string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("reference = ?", reference).Find(&transactions).Error
	return transactions, err
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
