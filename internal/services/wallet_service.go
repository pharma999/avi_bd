package services

import (
	"errors"
	"fmt"
	"time"

	"avi_bd/internal/dto"
	"avi_bd/internal/models"
	"avi_bd/internal/repositories"
	"avi_bd/internal/utils"
	"avi_bd/pkg/common"

	"gorm.io/gorm"
)

// WalletService handles wallet operations
type WalletService struct {
	walletRepo      *repositories.WalletRepository
	transactionRepo *repositories.TransactionRepository
}

// NewWalletService creates a new wallet service
func NewWalletService(walletRepo *repositories.WalletRepository, transactionRepo *repositories.TransactionRepository) *WalletService {
	return &WalletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

// GetWallet retrieves a user's wallet
func (s *WalletService) GetWallet(userID uint) (*dto.WalletResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	return &dto.WalletResponse{
		ID:          utils.UintToUUID(wallet.ID),
		UserID:      utils.UintToUUID(wallet.UserID),
		Balance:     wallet.Balance,
		Currency:    "USD", // Default currency
		LastUpdated: wallet.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetBalance retrieves user's balance
func (s *WalletService) GetBalance(userID uint) (float64, error) {
	return s.walletRepo.GetBalance(userID)
}

// Credit adds money to the wallet (transaction-safe within caller's transaction)
func (s *WalletService) Credit(userID uint, amount float64, txType string, reference string) error {
	if amount < 0 {
		return fmt.Errorf("credit amount cannot be negative")
	}

	// Increment balance
	if err := s.walletRepo.IncrementBalance(userID, amount); err != nil {
		return err
	}

	// Record transaction
	transaction := &struct {
		UserID    uint
		Type      string
		Amount    float64
		Reference string
	}{
		UserID:    userID,
		Type:      txType,
		Amount:    amount,
		Reference: reference,
	}

	_, _ = s.transactionRepo, transaction // suppress unused warning temporarily

	return nil
}

// Deposit adds money (simulated payment — no real gateway)
func (s *WalletService) Deposit(userID uint, amount float64) (*dto.WalletResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be greater than 0")
	}
	if err := s.walletRepo.IncrementBalance(userID, amount); err != nil {
		return nil, err
	}
	s.transactionRepo.Create(&models.Transaction{
		UserID:    userID,
		Type:      "deposit",
		Amount:    amount,
		Reference: fmt.Sprintf("deposit_sim_%d", time.Now().UnixNano()),
	})
	return s.GetWallet(userID)
}

// Withdraw deducts money (simulated — no real gateway)
func (s *WalletService) Withdraw(userID uint, amount float64) (*dto.WalletResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be greater than 0")
	}
	balance, err := s.walletRepo.GetBalance(userID)
	if err != nil {
		return nil, err
	}
	if balance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}
	if err := s.walletRepo.DecrementBalance(userID, amount); err != nil {
		return nil, err
	}
	s.transactionRepo.Create(&models.Transaction{
		UserID:    userID,
		Type:      "withdrawal",
		Amount:    amount,
		Reference: fmt.Sprintf("withdraw_sim_%d", time.Now().UnixNano()),
	})
	return s.GetWallet(userID)
}

// Debit removes money from the wallet (transaction-safe within caller's transaction)
func (s *WalletService) Debit(userID uint, amount float64, txType string, reference string) error {
	if amount < 0 {
		return fmt.Errorf("debit amount cannot be negative")
	}

	// Check balance
	balance, err := s.walletRepo.GetBalance(userID)
	if err != nil {
		return err
	}

	if balance < amount {
		return common.ErrInsufficientBalance
	}

	// Decrement balance
	if err := s.walletRepo.DecrementBalance(userID, amount); err != nil {
		return err
	}

	// Record transaction
	transaction := &models.Transaction{
		UserID:    userID,
		Type:      txType,
		Amount:    amount,
		Reference: reference,
	}
	s.transactionRepo.Create(transaction)

	return nil
}
