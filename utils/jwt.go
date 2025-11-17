package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/antoniocfetngnu/users-api/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"sub"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

var cfg *config.Config

func InitJWT(c *config.Config) {
	cfg = c
}

// GenerateJWT generates a new JWT token (for login)
func GenerateJWT(userID uint, username, email string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID), // CRITICAL: Kong uses this
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Add "kid" header for Kong compatibility
	token.Header["kid"] = "jwt-issuer-key"

	return token.SignedString([]byte(cfg.JWTSecret))
}

// DecodeJWTPayload extracts user info from JWT WITHOUT verifying signature
// (Kong already validated it, we just need to read the payload)
func DecodeJWTPayload(tokenString string) (userID uint, username, email string, err error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return 0, "", "", fmt.Errorf("invalid token format")
	}

	// Decode payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to decode payload: %w", err)
	}

	// Parse JSON
	var claims struct {
		Sub      json.Number `json:"sub"`
		Username string      `json:"username"`
		Email    string      `json:"email"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, "", "", fmt.Errorf("failed to parse payload: %w", err)
	}

	// Extract user ID from 'sub' claim
	subStr := claims.Sub.String()
	if subStr == "" {
		return 0, "", "", fmt.Errorf("missing sub claim")
	}

	var uid uint64
	if _, err := fmt.Sscanf(subStr, "%d", &uid); err != nil {
		return 0, "", "", fmt.Errorf("invalid sub claim: %w", err)
	}

	return uint(uid), claims.Username, claims.Email, nil
}
