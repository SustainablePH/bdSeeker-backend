package middleware

import (
	"net/http"
	"strings"

	"github.com/bishworup11/bdSeeker-backend/internal/config"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role"
)

// AuthMiddleware validates JWT tokens from cookies or Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First, try to get token from cookie (preferred for browser clients)
		token, err := utils.GetAuthCookie(c.Request)
		if err == nil && token != "" {
			tokenString = token
		} else {
			// Fallback to Authorization header (for API clients like Postman)
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				c.Abort()
				return
			}
			tokenString = parts[1]
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			// Clear cookie if token is invalid
			utils.ClearAuthCookie(c.Writer)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in Gin context
		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(UserEmailKey), claims.Email)
		c.Set(string(UserRoleKey), claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if the user has one of the allowed roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleVal, exists := c.Get(string(UserRoleKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		userRole, ok := userRoleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user role format"})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID extracts user ID from Gin context
func GetUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get(string(UserIDKey))
	if !exists {
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	return userID, ok
}

// GetUserEmail extracts user email from Gin context
func GetUserEmail(c *gin.Context) (string, bool) {
	emailVal, exists := c.Get(string(UserEmailKey))
	if !exists {
		return "", false
	}
	email, ok := emailVal.(string)
	return email, ok
}

// GetUserRole extracts user role from Gin context
func GetUserRole(c *gin.Context) (string, bool) {
	roleVal, exists := c.Get(string(UserRoleKey))
	if !exists {
		return "", false
	}
	role, ok := roleVal.(string)
	return role, ok
}
