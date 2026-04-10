package dto

// TransactionItem represents a single transaction per specification
type TransactionItem struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Type          string  `json:"type"` // bet, cashout, win, withdrawal, deposit, loss
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	Reference     string  `json:"reference"`
	Description   string  `json:"description"`
	Status        string  `json:"status"` // completed, pending, failed
	Timestamp     string  `json:"timestamp"`
}

// TransactionHistoryResponse represents transaction history response per specification
type TransactionHistoryResponse struct {
	Transactions []TransactionItem `json:"transactions"`
}

// TransactionRequest represents a transaction request per specification
type TransactionRequest struct {
	UserID    string  `json:"user_id"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

// TransactionHistoryListResponse wraps paginated transactions response
type TransactionHistoryListResponse struct {
	Transactions []TransactionItem `json:"transactions"`
}
