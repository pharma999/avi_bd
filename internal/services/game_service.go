package services

import (
	"errors"

	"avi_bd/internal/repositories"
	"avi_bd/pkg/common"

	"gorm.io/gorm"
)

// GameService handles game round operations
type GameService struct {
	roundRepo *repositories.GameRoundRepository
	betRepo   *repositories.BetRepository
}

// NewGameService creates a new game service
func NewGameService(roundRepo *repositories.GameRoundRepository, betRepo *repositories.BetRepository) *GameService {
	return &GameService{
		roundRepo: roundRepo,
		betRepo:   betRepo,
	}
}

// GetGameHistory retrieves game round history
func (s *GameService) GetGameHistory(limit int, offset int) (interface{}, int64, error) {
	rounds, count, err := s.roundRepo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return rounds, count, nil
}

// GetCurrentRound retrieves the current game round
func (s *GameService) GetCurrentRound() (interface{}, error) {
	round, err := s.roundRepo.GetLatestRound()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRoundNotFound
		}
		return nil, err
	}
	return round, nil
}
