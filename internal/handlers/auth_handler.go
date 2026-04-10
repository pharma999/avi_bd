package handlers

import (
	"net/http"

	"avi_bd/internal/dto"
	"avi_bd/internal/services"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with name, email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.AuthResponse "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input or validation failed"
// @Failure 409 {object} map[string]interface{} "Email already exists"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Register user
	response, err := h.authService.Register(&req)
	if err != nil {
		// Check if it's a "user exists" error
		if err.Error() == "user already exists" {
			utils.ErrorResponse(c, http.StatusConflict, "Email already exists", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Registration failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

// Login handles user login
// @Summary User login
// @Description Login with email and password to get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.AuthResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid credentials or validation failed"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Login user
	response, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// Verify handles token verification
// @Summary Verify JWT token
// @Description Verify the validity of a JWT token and return user information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.VerifyTokenResponse "Token is valid"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /verify [get]
func (h *AuthHandler) Verify(c *gin.Context) {
	// Get user ID from context (set by middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	// Convert userID to uint
	uid, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid user ID type")
		return
	}

	// Verify token and get user info
	response, err := h.authService.VerifyToken(uid)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Verification failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token verified successfully", response)
}
