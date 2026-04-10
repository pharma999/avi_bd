package handlers

import (
	"net/http"

	"avi_bd/internal/dto"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Health checks the health of the server
// @Summary Health check
// @Description Get server health status including database connection
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse "Server is healthy"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	healthResp := &dto.HealthCheckResponse{
		Status:    "healthy",
		Timestamp: utils.GetTimestamp(),
		Database:  "connected",
	}

	// Check database connection
	if h.db == nil {
		healthResp.Database = "disconnected"
		healthResp.Status = "degraded"
	} else if err := h.db.Exec("SELECT 1").Error; err != nil {
		healthResp.Database = "disconnected"
		healthResp.Status = "degraded"
	}

	utils.SuccessResponse(c, http.StatusOK, "Server health check", healthResp)
}
