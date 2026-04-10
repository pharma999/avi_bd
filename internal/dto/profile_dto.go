package dto

// ProfileResponse represents user profile response per specification
type ProfileResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Balance   float64 `json:"balance"`
	TotalBets int     `json:"total_bets"`
	TotalWins int     `json:"total_wins"`
	WinRate   float64 `json:"win_rate"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// UpdateProfileRequest represents profile update request per specification
type UpdateProfileRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=100"`
	Email *string `json:"email" binding:"omitempty,email"`
}
