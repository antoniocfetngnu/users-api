package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/antoniocfetngnu/users-api/utils"
	"github.com/gin-gonic/gin"
)

// JWTPayload represents the decoded JWT payload
type JWTPayload struct {
	Sub      json.Number `json:"sub"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
}

// AuthMiddleware checks for Kong-validated JWT, then falls back to cookie
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Priority 1: Check if Kong validated the JWT (cookie exists and request came through Kong)
		cookie, err := c.Cookie("auth_token")
		if err == nil && cookie != "" {
			// Kong already validated this JWT (signature, expiration, etc.)
			// We just need to extract the user ID from the 'sub' claim

			// Decode JWT payload (no need to verify signature - Kong did that)
			parts := strings.Split(cookie, ".")
			if len(parts) == 3 {
				// Decode payload (second part)
				payload, err := base64.RawURLEncoding.DecodeString(parts[1])
				if err == nil {
					var jwtPayload JWTPayload
					if json.Unmarshal(payload, &jwtPayload) == nil {
						// Extract user ID from 'sub' claim
						if subStr := jwtPayload.Sub.String(); subStr != "" {
							userID, err := strconv.ParseUint(subStr, 10, 32)
							if err == nil {
								// Successfully extracted user ID
								c.Set("userID", uint(userID))
								c.Set("username", jwtPayload.Username)
								c.Set("email", jwtPayload.Email)
								c.Set("authSource", "jwt-decoded")
								c.Next()
								return
							}
						}
					}
				}
			}

			// If decoding failed, fall back to full JWT validation
			claims, err := utils.ValidateJWT(cookie)
			if err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("email", claims.Email)
				c.Set("authSource", "jwt-validated")
				c.Next()
				return
			}
		}

		// Unauthorized if we get here
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
