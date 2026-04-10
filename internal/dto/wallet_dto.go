package dto

// WalletResponse represents wallet data in response per specification
type WalletResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	LastUpdated string  `json:"last_updated"`
}
