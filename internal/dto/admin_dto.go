package dto

// AdminStatsResponse represents platform-wide statistics
type AdminStatsResponse struct {
	TotalUsers   int64   `json:"total_users"`
	TotalGames   int64   `json:"total_games"`
	TotalBets    int64   `json:"total_bets"`
	TotalWagered float64 `json:"total_wagered"`
	TotalPayouts float64 `json:"total_payouts"`
	HouseRevenue float64 `json:"house_revenue"`
}

// AdminUserItem represents a single user entry in the admin user list
type AdminUserItem struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Balance       float64 `json:"balance"`
	TotalBets     int64   `json:"total_bets"`
	TotalWagered  float64 `json:"total_wagered"`
	TotalWinnings float64 `json:"total_winnings"`
	CreatedAt     string  `json:"created_at"`
}

// AdminUsersResponse wraps the user list
type AdminUsersResponse struct {
	Users []AdminUserItem `json:"users"`
}
