package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/pkg/utils"
)

const ( // Define keys for storing user info in Gin context
	ContextUserIDKey = "userID"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(jwtCfg configs.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Expect "Bearer TOKEN"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString, jwtCfg.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store UserID in context
		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserIDFromContext retrieves the UserID from the Gin context.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, ok := c.Get(ContextUserIDKey)
	if !ok {
		return uuid.Nil, false
	}

	uuidUserID, ok := userID.(uuid.UUID)
	return uuidUserID, ok
}
