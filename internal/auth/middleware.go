package auth

import (
	"net/http"
	"strings"

	"avi_bd/config"
	"avi_bd/internal/utils"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a middleware that validates JWT tokens
func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header", nil)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := ValidateToken(tokenString, cfg)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserID retrieves the user ID from context
func GetUserID(c *gin.Context) uint {
	userID, _ := c.Get("userID")
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// GetEmail retrieves the email from context
func GetEmail(c *gin.Context) string {
	email, _ := c.Get("email")
	if e, ok := email.(string); ok {
		return e
	}
	return ""
}
