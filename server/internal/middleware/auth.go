package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
	"zhanbu/pkg/utils"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens.
// It extracts the user ID and username from the token and sets them in the context.
func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized,
				apperrors.GetMessage(apperrors.ErrUnauthorized))
			c.Abort()
			return
		}

		// Expect "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized,
				"invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrTokenInvalid,
				apperrors.GetMessage(apperrors.ErrTokenInvalid))
			c.Abort()
			return
		}

		// Set user info in context for handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// OptionalAuthMiddleware extracts user info when a valid Bearer token is present.
// Public endpoints can use it to persist data for logged-in users without
// rejecting anonymous requests.
func OptionalAuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}
		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context.
// Returns 0 if not found.
func GetUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		if uid, ok := id.(uint); ok {
			return uid
		}
	}
	return 0
}

// GetUsername extracts the username from the Gin context.
func GetUsername(c *gin.Context) string {
	if name, exists := c.Get("username"); exists {
		if uname, ok := name.(string); ok {
			return uname
		}
	}
	return ""
}
