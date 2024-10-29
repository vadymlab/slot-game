package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

// GenerateToken creates a signed JWT token for a given user ID with a specified lifetime.
// The token includes standard claims, such as expiration time, issue time, user ID (as the subject), and a unique token ID.
// Returns the signed token string or an error if signing fails.
func GenerateToken(userID *uuid.UUID, secret string, lifeTime int) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(lifeTime) * time.Minute)), // Token expiration time
		IssuedAt:  jwt.NewNumericDate(time.Now()),                                            // Token issue time
		Subject:   userID.String(),                                                           // User ID as the subject
		ID:        uuid.NewString(),                                                          // Unique token ID
	}

	// Create a new token with HS256 signing method and add claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret)) // Sign the token with the provided secret
}
