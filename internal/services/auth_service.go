package services

import (
	"errors"
	"time"

	"avi_bd/config"
	"avi_bd/internal/auth"
	"avi_bd/internal/dto"
	"avi_bd/internal/models"
	"avi_bd/internal/repositories"
	"avi_bd/internal/utils"
	"avi_bd/pkg/common"

	"gorm.io/gorm"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo   *repositories.UserRepository
	walletRepo *repositories.WalletRepository
	cfg        *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repositories.UserRepository, walletRepo *repositories.WalletRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		cfg:        cfg,
	}
}

// Register registers a new user
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create wallet for user
	wallet := &models.Wallet{
		UserID:  user.ID,
		Balance: s.cfg.InitialBalance,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        utils.UintToUUID(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Balance:   wallet.Balance,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// Login logs in a user
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, common.ErrInvalidCredentials
	}

	// Get wallet for balance
	wallet, err := s.walletRepo.GetByUserID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        utils.UintToUUID(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Balance:   wallet.Balance,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// VerifyToken verifies a JWT token and returns user info with balance
func (s *AuthService) VerifyToken(userID uint) (*dto.VerifyTokenResponse, error) {
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

	return &dto.VerifyTokenResponse{
		ID:        utils.UintToUUID(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Balance:   wallet.Balance,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
