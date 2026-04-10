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

// BetHandler handles bet endpoints
type BetHandler struct {
	betService *services.BetService
}

// NewBetHandler creates a new bet handler
func NewBetHandler(betService *services.BetService) *BetHandler {
	return &BetHandler{
		betService: betService,
	}
}

// PlaceBet places a bet
// @Summary Place a new bet
// @Description Place a bet on the current game round
// @Tags Betting
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BetRequest true "Bet placement data"
// @Success 201 {object} dto.BetResponse "Bet placed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input or insufficient balance"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Validation failed"
// @Router /bet [post]
func (h *BetHandler) PlaceBet(c *gin.Context) {
	userID := auth.GetUserID(c)
	var req dto.BetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Validation failed", errors)
		return
	}

	// Place bet
	response, err := h.betService.PlaceBet(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Bet placed successfully", response)
}

// Cashout cashes out a bet
// @Summary Cashout an active bet
// @Description Cashout (stop) an active bet at current multiplier
// @Tags Betting
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CashoutRequest true "Cashout request data"
// @Success 200 {object} dto.CashoutResponse "Cashout successful"
// @Failure 400 {object} map[string]interface{} "Invalid multiplier"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Bet not found"
// @Failure 422 {object} map[string]interface{} "Bet already closed"
// @Router /cashout [post]
func (h *BetHandler) Cashout(c *gin.Context) {
	userID := auth.GetUserID(c)
	var req dto.CashoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate multiplier
	if req.CurrentMultiplier < 1.0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid multiplier", nil)
		return
	}

	// Convert BetID from UUID string to uint
	betID, err := utils.UUIDStringToUint(req.BetID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid bet ID format", nil)
		return
	}

	// Cashout bet
	response, err := h.betService.Cashout(userID, betID, req.CurrentMultiplier)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cashout successful", response)
}

// GetBetHistory retrieves user's bet history
// @Summary Get bet history
// @Description Get user's bet history with pagination (max 100 per page)
// @Tags Betting
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20) maximum(100)
// @Param status query string false "Filter by status"
// @Success 200 {object} object "Bet history with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Validation failed"
// @Router /bets/history [get]
func (h *BetHandler) GetBetHistory(c *gin.Context) {
	userID := auth.GetUserID(c)

	// Parse and validate pagination
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			if parsed > 100 {
				utils.ErrorResponse(c, http.StatusUnprocessableEntity, "Limit cannot exceed 100", nil)
				return
			}
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	// Get bets
	bets, total, err := h.betService.GetBetHistory(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve bet history", err.Error())
		return
	}

	// Build response
	response := &dto.BetHistoryResponse{
		Bets: bets,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, response, page, limit, total)
}
