package services

import (
	"errors"
	"time"

	"avi_bd/internal/dto"
	"avi_bd/internal/repositories"
	"avi_bd/internal/utils"
	"avi_bd/pkg/common"

	"gorm.io/gorm"
)

// ProfileService handles user profile operations
type ProfileService struct {
	userRepo   *repositories.UserRepository
	walletRepo *repositories.WalletRepository
}

// NewProfileService creates a new profile service
func NewProfileService(userRepo *repositories.UserRepository, walletRepo *repositories.WalletRepository) *ProfileService {
	return &ProfileService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

// GetProfile retrieves user's profile
func (s *ProfileService) GetProfile(userID uint) (*dto.ProfileResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	// Calculate statistics (placeholder for now - can be enhanced later with actual bet queries)
	totalBets := 0
	totalWins := 0
	winRate := 0.0

	return &dto.ProfileResponse{
		ID:        utils.UintToUUID(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Balance:   wallet.Balance,
		TotalBets: totalBets,
		TotalWins: totalWins,
		WinRate:   winRate,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateProfile updates user's profile
func (s *ProfileService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	// Update name if provided
	if req.Name != nil {
		user.Name = *req.Name
	}

	// Update email if provided
	if req.Email != nil {
		// Check if email is taken by another user
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	// Save updates
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Get updated wallet
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:        utils.UintToUUID(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Balance:   wallet.Balance,
		TotalBets: 0,
		TotalWins: 0,
		WinRate:   0,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetStatistics retrieves user betting statistics
func (s *ProfileService) GetStatistics(userID uint) (*dto.StatisticsResponse, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	// Return statistics (placeholder for now)
	return &dto.StatisticsResponse{
		TotalBets:         0,
		TotalWins:         0,
		TotalLosses:       0,
		TotalWagered:      0,
		TotalWinnings:     0,
		AverageMultiplier: 0,
		WinRate:           0,
		HighestMultiplier: 0,
		LargestWin:        0,
		LargestLoss:       0,
	}, nil
}
