package middleware

import (
	"github.com/antoniocfetngnu/users-api/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware extracts user info from JWT (Kong already validated it)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT from cookie (Kong already validated this)
		cookie, err := c.Cookie("auth_token")
		if err != nil || cookie == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract user info from JWT payload (no signature verification needed)
		userID, username, email, err := utils.DecodeJWTPayload(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", userID)
		c.Set("username", username)
		c.Set("email", email)
		c.Next()
	}
}

// OptionalAuthMiddleware tries to extract user info but doesn't fail if missing
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth_token")
		if err == nil && cookie != "" {
			userID, username, email, err := utils.DecodeJWTPayload(cookie)
			if err == nil {
				c.Set("userID", userID)
				c.Set("username", username)
				c.Set("email", email)
			}
		}
		c.Next() // Continue even if not authenticated
	}
}
