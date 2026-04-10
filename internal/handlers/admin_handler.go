package handlers

import (
	"net/http"
	"strconv"

	"avi_bd/internal/dto"
	"avi_bd/internal/game"
	"avi_bd/internal/services"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin/analytics endpoints
type AdminHandler struct {
	adminService *services.AdminService
	engine       *game.Engine
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminService *services.AdminService, engine *game.Engine) *AdminHandler {
	return &AdminHandler{adminService: adminService, engine: engine}
}

// GetStats returns platform-wide statistics
// GET /api/v1/admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve stats", nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Stats retrieved", stats)
}

// GetUsers returns all users with aggregated bet stats
// GET /api/v1/admin/users?page=1&limit=50
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, limit := parsePagination(c, 50)
	offset := (page - 1) * limit

	users, total, err := h.adminService.GetAllUsers(limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", nil)
		return
	}

	response := &dto.AdminUsersResponse{Users: users}
	utils.PaginatedSuccessResponse(c, http.StatusOK, response, page, limit, total)
}

// GetUserBets returns paginated bets for a specific user
// GET /api/v1/admin/users/:id/bets?page=1&limit=20
func (h *AdminHandler) GetUserBets(c *gin.Context) {
	rawID := c.Param("id")
	userID, err := utils.UUIDStringToUint(rawID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	page, limit := parsePagination(c, 20)
	offset := (page - 1) * limit

	bets, total, err := h.adminService.GetUserBets(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user bets", nil)
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, gin.H{"bets": bets}, page, limit, total)
}

// GetCrashMode returns current crash mode (auto / manual)
// GET /api/v1/admin/crash-mode
func (h *AdminHandler) GetCrashMode(c *gin.Context) {
	mode := "auto"
	if h.engine.IsManualMode() {
		mode = "manual"
	}
	utils.SuccessResponse(c, http.StatusOK, "OK", gin.H{"mode": mode})
}

// SetCrashMode switches between auto and manual crash mode
// POST /api/v1/admin/crash-mode   body: {"mode":"manual"} or {"mode":"auto"}
func (h *AdminHandler) SetCrashMode(c *gin.Context) {
	var body struct {
		Mode string `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "mode field required (auto|manual)", nil)
		return
	}
	switch body.Mode {
	case "manual":
		h.engine.SetManualMode(true)
	case "auto":
		h.engine.SetManualMode(false)
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "mode must be 'auto' or 'manual'", nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Crash mode updated", gin.H{"mode": body.Mode})
}

// ForceCrash triggers an immediate crash (manual mode only)
// POST /api/v1/admin/force-crash
func (h *AdminHandler) ForceCrash(c *gin.Context) {
	if !h.engine.IsManualMode() {
		utils.ErrorResponse(c, http.StatusBadRequest, "Engine is in auto mode — switch to manual first", nil)
		return
	}
	ok := h.engine.ForceCrash()
	if !ok {
		utils.ErrorResponse(c, http.StatusConflict, "No round is currently flying", nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Crash triggered", nil)
}

// parsePagination extracts page/limit from query params with a default limit
func parsePagination(c *gin.Context, defaultLimit int) (page, limit int) {
	page = 1
	limit = defaultLimit
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 200 {
			limit = v
		}
	}
	return
}
