package repositories

import (
	"avi_bd/internal/models"

	"gorm.io/gorm"
)

// WalletRepository handles wallet database operations
type WalletRepository struct {
	db *gorm.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create creates a new wallet
func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

// GetByUserID retrieves a wallet by user ID
func (r *WalletRepository) GetByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Update updates a wallet
func (r *WalletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

// IncrementBalance increments wallet balance (for atomic operations)
func (r *WalletRepository) IncrementBalance(userID uint, amount float64) error {
	return r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// DecrementBalance decrements wallet balance (for atomic operations)
func (r *WalletRepository) DecrementBalance(userID uint, amount float64) error {
	return r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance - ?", amount)).Error
}

// GetBalance retrieves current balance
func (r *WalletRepository) GetBalance(userID uint) (float64, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}
