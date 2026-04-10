package handlers

import (
	"net/http"

	"avi_bd/internal/auth"
	"avi_bd/internal/services"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
)

// WalletHandler handles wallet endpoints
type WalletHandler struct {
	walletService *services.WalletService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// GetWallet retrieves user's wallet
// @Summary Get wallet balance
// @Description Get user's current wallet balance and details
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.WalletResponse "Wallet retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Wallet retrieval failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /wallet [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := auth.GetUserID(c)
	wallet, err := h.walletService.GetWallet(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Wallet not found", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Wallet retrieved", wallet)
}

// Deposit adds balance (simulated payment — no real money)
// POST /api/v1/wallet/deposit   body: {"amount": 500}
func (h *WalletHandler) Deposit(c *gin.Context) {
	userID := auth.GetUserID(c)
	var body struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "amount must be greater than 0", nil)
		return
	}
	wallet, err := h.walletService.Deposit(userID, body.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Deposit successful", wallet)
}

// Withdraw deducts balance (simulated — no real money)
// POST /api/v1/wallet/withdraw   body: {"amount": 200}
func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID := auth.GetUserID(c)
	var body struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "amount must be greater than 0", nil)
		return
	}
	wallet, err := h.walletService.Withdraw(userID, body.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Withdrawal successful", wallet)
}
