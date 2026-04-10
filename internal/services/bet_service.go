package services

import (
	"errors"
	"fmt"
	"time"

	"avi_bd/internal/constants"
	"avi_bd/internal/dto"
	"avi_bd/internal/models"
	"avi_bd/internal/repositories"
	"avi_bd/internal/utils"
	"avi_bd/pkg/common"

	"gorm.io/gorm"
)

// BetService handles betting operations
type BetService struct {
	betRepo         *repositories.BetRepository
	roundRepo       *repositories.GameRoundRepository
	walletRepo      *repositories.WalletRepository
	transactionRepo *repositories.TransactionRepository
}

// NewBetService creates a new bet service
func NewBetService(
	betRepo *repositories.BetRepository,
	roundRepo *repositories.GameRoundRepository,
	walletRepo *repositories.WalletRepository,
	transactionRepo *repositories.TransactionRepository,
) *BetService {
	return &BetService{
		betRepo:         betRepo,
		roundRepo:       roundRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

// PlaceBet places a bet on the current round
func (s *BetService) PlaceBet(userID uint, req *dto.BetRequest) (*dto.BetResponse, error) {
	// Validate bet amount
	if req.Amount <= 0 {
		return nil, common.ErrInvalidBetAmount
	}

	// Get current round
	round, err := s.roundRepo.GetLatestRound()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRoundNotFound
		}
		return nil, err
	}

	// Check if round is in waiting state
	if round.Status != constants.GameStatusWaiting {
		return nil, errors.New("betting is only allowed during waiting state")
	}

	// Check wallet balance
	balance, err := s.walletRepo.GetBalance(userID)
	if err != nil {
		return nil, err
	}

	if balance < req.Amount {
		return nil, common.ErrInsufficientBalance
	}

	// Debit wallet
	if err := s.walletRepo.DecrementBalance(userID, req.Amount); err != nil {
		return nil, err
	}

	// Create bet
	bet := &models.Bet{
		UserID:      userID,
		RoundID:     round.ID,
		Amount:      req.Amount,
		AutoCashout: req.AutoCashout,
		Result:      constants.BetResultPending,
		Payout:      0,
	}

	if err := s.betRepo.Create(bet); err != nil {
		// Rollback wallet if bet creation fails
		s.walletRepo.IncrementBalance(userID, req.Amount)
		return nil, err
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:    userID,
		Type:      constants.TransactionTypeBet,
		Amount:    req.Amount,
		Reference: fmt.Sprintf("bet_%d", bet.ID),
	}
	s.transactionRepo.Create(transaction)

	return &dto.BetResponse{
		ID:          utils.UintToUUID(bet.ID),
		UserID:      utils.UintToUUID(bet.UserID),
		RoundID:     utils.UintToUUID(bet.RoundID),
		Amount:      bet.Amount,
		AutoCashout: bet.AutoCashout,
		Status:      constants.BetStatusActive,
		PlacedAt:    bet.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Cashout cashes out a bet
func (s *BetService) Cashout(userID uint, betID uint, currentMultiplier float64) (*dto.CashoutResponse, error) {
	// Get bet
	bet, err := s.betRepo.GetByID(betID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrBetNotFound
		}
		return nil, err
	}

	// Verify bet belongs to user
	if bet.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if already cashed out or lost
	if bet.Result != constants.BetResultPending {
		return nil, common.ErrAlreadyCashedOut
	}

	// Get round
	round, err := s.roundRepo.GetByID(bet.RoundID)
	if err != nil {
		return nil, err
	}

	// Check if round has crashed
	if round.Status == constants.GameStatusCrashed {
		return nil, common.ErrRoundCrashed
	}

	// Calculate payout
	payout := bet.Amount * currentMultiplier

	// Update bet
	bet.Result = constants.BetResultWin
	bet.Payout = payout
	bet.CashoutMultiplier = currentMultiplier

	if err := s.betRepo.Update(bet); err != nil {
		return nil, err
	}

	// Credit wallet
	if err := s.walletRepo.IncrementBalance(userID, payout); err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:    userID,
		Type:      constants.TransactionTypeCashout,
		Amount:    payout,
		Reference: fmt.Sprintf("bet_%d", bet.ID),
	}
	s.transactionRepo.Create(transaction)

	return &dto.CashoutResponse{
		ID:                utils.UintToUUID(bet.ID),
		UserID:            utils.UintToUUID(bet.UserID),
		Status:            constants.BetResultCashedOut,
		CashoutMultiplier: currentMultiplier,
		Winnings:          payout,
		CashedOutAt:       bet.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetBetHistory retrieves user's bet history
func (s *BetService) GetBetHistory(userID uint, limit int, offset int) ([]dto.BetHistoryItem, int64, error) {
	bets, count, err := s.betRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var response []dto.BetHistoryItem
	for _, bet := range bets {
		// Get round details
		round, _ := s.roundRepo.GetByID(bet.RoundID)

		historyBet := dto.BetHistoryItem{
			ID:       utils.UintToUUID(bet.ID),
			RoundID:  utils.UintToUUID(bet.RoundID),
			Amount:   bet.Amount,
			Status:   bet.Result,
			PlacedAt: bet.CreatedAt.Format(time.RFC3339),
		}

		// Only include optional fields if they have values
		if bet.Result != constants.BetResultPending {
			historyBet.Winnings = utils.Float64Ptr(bet.Payout)
			if round != nil {
				historyBet.FinalMultiplier = utils.Float64Ptr(round.CrashMultiplier)
			}
			if bet.CashoutMultiplier > 0 {
				historyBet.CashoutMultiplier = utils.Float64Ptr(bet.CashoutMultiplier)
			}
			endedAt := bet.UpdatedAt.Format(time.RFC3339)
			historyBet.EndedAt = &endedAt
		}

		response = append(response, historyBet)
	}

	return response, count, nil
}
