package services

import (
	"avi_bd/internal/dto"
	"avi_bd/internal/models"
	"avi_bd/internal/utils"

	"gorm.io/gorm"
)

// AdminService handles admin/analytics operations
type AdminService struct {
	db *gorm.DB
}

// NewAdminService creates a new admin service
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// GetStats returns platform-wide statistics
func (s *AdminService) GetStats() (*dto.AdminStatsResponse, error) {
	var totalUsers int64
	s.db.Model(&models.User{}).Count(&totalUsers)

	var totalGames int64
	s.db.Model(&models.GameRound{}).Count(&totalGames)

	var totalBets int64
	s.db.Model(&models.Bet{}).Count(&totalBets)

	var totalWagered float64
	s.db.Model(&models.Bet{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalWagered)

	var totalPayouts float64
	s.db.Model(&models.Bet{}).
		Where("result IN ('won', 'cashed_out')").
		Select("COALESCE(SUM(payout), 0)").
		Scan(&totalPayouts)

	return &dto.AdminStatsResponse{
		TotalUsers:   totalUsers,
		TotalGames:   totalGames,
		TotalBets:    totalBets,
		TotalWagered: totalWagered,
		TotalPayouts: totalPayouts,
		HouseRevenue: totalWagered - totalPayouts,
	}, nil
}

// GetAllUsers returns all users with aggregated stats, paginated
func (s *AdminService) GetAllUsers(limit, offset int) ([]dto.AdminUserItem, int64, error) {
	var users []models.User
	var total int64

	s.db.Model(&models.User{}).Count(&total)

	if err := s.db.
		Preload("Wallet").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	result := make([]dto.AdminUserItem, 0, len(users))
	for _, u := range users {
		var betCount int64
		s.db.Model(&models.Bet{}).Where("user_id = ?", u.ID).Count(&betCount)

		var totalWagered float64
		s.db.Model(&models.Bet{}).
			Where("user_id = ?", u.ID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalWagered)

		var totalWinnings float64
		s.db.Model(&models.Bet{}).
			Where("user_id = ? AND result IN ('won', 'cashed_out')", u.ID).
			Select("COALESCE(SUM(payout), 0)").
			Scan(&totalWinnings)

		balance := 0.0
		if u.Wallet != nil {
			balance = u.Wallet.Balance
		}

		result = append(result, dto.AdminUserItem{
			ID:            utils.UintToUUID(u.ID),
			Username:      u.Name,
			Email:         u.Email,
			Balance:       balance,
			TotalBets:     betCount,
			TotalWagered:  totalWagered,
			TotalWinnings: totalWinnings,
			CreatedAt:     u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, total, nil
}

// GetUserBets returns paginated bets for a specific user
func (s *AdminService) GetUserBets(userID uint, limit, offset int) ([]dto.BetHistoryItem, int64, error) {
	var bets []models.Bet
	var total int64

	s.db.Model(&models.Bet{}).Where("user_id = ?", userID).Count(&total)

	if err := s.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&bets).Error; err != nil {
		return nil, 0, err
	}

	result := make([]dto.BetHistoryItem, 0, len(bets))
	for _, bet := range bets {
		var round models.GameRound
		s.db.First(&round, bet.RoundID)

		item := dto.BetHistoryItem{
			ID:       utils.UintToUUID(bet.ID),
			RoundID:  utils.UintToUUID(bet.RoundID),
			Amount:   bet.Amount,
			Status:   bet.Result,
			PlacedAt: bet.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if bet.Result != "pending" {
			item.Winnings = utils.Float64Ptr(bet.Payout)
			item.FinalMultiplier = utils.Float64Ptr(round.CrashMultiplier)
			if bet.CashoutMultiplier > 0 {
				item.CashoutMultiplier = utils.Float64Ptr(bet.CashoutMultiplier)
			}
			endedAt := bet.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
			item.EndedAt = &endedAt
		}
		result = append(result, item)
	}

	return result, total, nil
}
