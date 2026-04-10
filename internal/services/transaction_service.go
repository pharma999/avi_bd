package services

import (
	"avi_bd/internal/repositories"
)

// TransactionService handles transaction operations
type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// GetTransactionHistory retrieves user's transaction history
func (s *TransactionService) GetTransactionHistory(userID uint, limit int, offset int) (interface{}, int64, error) {
	transactions, count, err := s.transactionRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return transactions, count, nil
}
