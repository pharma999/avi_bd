package handlers

import (
	"net/http"

	"avi_bd/internal/auth"
	"avi_bd/internal/dto"
	"avi_bd/internal/services"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles profile endpoints
type ProfileHandler struct {
	profileService *services.ProfileService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile retrieves user's profile
// @Summary Get user profile
// @Description Get user's profile including name, email, and wallet balance
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ProfileResponse "Profile retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Profile retrieval failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := auth.GetUserID(c)

	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Profile not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved", profile)
}

// UpdateProfile updates user's profile
// @Summary Update user profile
// @Description Update user's name and/or email
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Updated profile data"
// @Success 200 {object} dto.ProfileResponse "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Email already taken"
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := auth.GetUserID(c)
	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Update profile
	profile, err := h.profileService.UpdateProfile(userID, &req)
	if err != nil {
		// Check if it's email conflict
		if err.Error() == "email already exists" {
			utils.ErrorResponse(c, http.StatusConflict, "Email already taken", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Profile update failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", profile)
}

// GetStatistics retrieves user betting statistics
// @Summary Get user statistics
// @Description Get user's betting statistics including total bets, wins, losses, and performance metrics
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.StatisticsResponse "Statistics retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /profile/statistics [get]
func (h *ProfileHandler) GetStatistics(c *gin.Context) {
	userID := auth.GetUserID(c)

	stats, err := h.profileService.GetStatistics(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Statistics not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved", stats)
}
