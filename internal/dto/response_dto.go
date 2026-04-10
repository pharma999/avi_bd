package dto

// VerifyTokenResponse represents token verification response per specification
type VerifyTokenResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

// StatisticsResponse represents user statistics per specification
type StatisticsResponse struct {
	TotalBets         int     `json:"total_bets"`
	TotalWins         int     `json:"total_wins"`
	TotalLosses       int     `json:"total_losses"`
	TotalWagered      float64 `json:"total_wagered"`
	TotalWinnings     float64 `json:"total_winnings"`
	AverageMultiplier float64 `json:"average_multiplier"`
	WinRate           float64 `json:"win_rate"`
	HighestMultiplier float64 `json:"highest_multiplier"`
	LargestWin        float64 `json:"largest_win"`
	LargestLoss       float64 `json:"largest_loss"`
}

// HealthCheckResponse represents health check response per specification
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database"`
}
