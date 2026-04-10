package handlers

import (
	"net/http"
	"strconv"

	"avi_bd/internal/auth"
	"avi_bd/internal/dto"
	"avi_bd/internal/services"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
)

// HistoryHandler handles history endpoints
type HistoryHandler struct {
	gameService        *services.GameService
	transactionService *services.TransactionService
}

// NewHistoryHandler creates a new history handler
func NewHistoryHandler(gameService *services.GameService, transactionService *services.TransactionService) *HistoryHandler {
	return &HistoryHandler{
		gameService:        gameService,
		transactionService: transactionService,
	}
}

// validatePagination validates pagination parameters
func validatePagination(c *gin.Context, defaultLimit int) (page int, limit int, valid bool) {
	page = 1
	limit = defaultLimit

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			if parsed > 100 {
				return page, limit, false // Invalid: exceeds max limit
			}
			limit = parsed
		}
	}

	return page, limit, true
}

// GetGameHistory retrieves game history
// @Summary Get game history
// @Description Get global game history with pagination (max 100 per page)
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(50) maximum(100)
// @Success 200 {object} dto.GameHistoryResponse "Game history retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Validation failed"
// @Router /game/history [get]
func (h *HistoryHandler) GetGameHistory(c *gin.Context) {
	page, limit, valid := validatePagination(c, 50)

	if !valid {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Limit cannot exceed 100", nil)
		return
	}

	offset := (page - 1) * limit

	history, total, err := h.gameService.GetGameHistory(limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve game history", err.Error())
		return
	}

	// Convert to proper response format
	rounds := make([]dto.GameRoundItem, 0)
	if historyList, ok := history.([]interface{}); ok {
		for _, item := range historyList {
			if round, ok := item.(dto.GameRoundItem); ok {
				rounds = append(rounds, round)
			}
		}
	}

	response := &dto.GameHistoryResponse{
		Rounds: rounds,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, response, page, limit, total)
}

// GetTransactionHistory retrieves user's transaction history
// @Summary Get transaction history
// @Description Get user's transaction history with pagination (max 100 per page)
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10) maximum(100)
// @Param type query string false "Transaction type filter"
// @Success 200 {object} dto.TransactionHistoryResponse "Transaction history retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Validation failed"
// @Router /transactions/history [get]
func (h *HistoryHandler) GetTransactionHistory(c *gin.Context) {
	userID := auth.GetUserID(c)
	page, limit, valid := validatePagination(c, 10)

	if !valid {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Limit cannot exceed 100", nil)
		return
	}

	offset := (page - 1) * limit

	transactions, total, err := h.transactionService.GetTransactionHistory(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transaction history", err.Error())
		return
	}

	// Convert to proper response format
	txItems := make([]dto.TransactionItem, 0)
	if txList, ok := transactions.([]interface{}); ok {
		for _, item := range txList {
			if tx, ok := item.(dto.TransactionItem); ok {
				txItems = append(txItems, tx)
			}
		}
	}

	response := &dto.TransactionHistoryResponse{
		Transactions: txItems,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, response, page, limit, total)
}
