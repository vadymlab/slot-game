package jwt

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vadymlab/slot-game/internal/constants"
	"net/http"
)

// AuthMiddleware is a middleware function for Gin that authenticates requests using a JWT token.
// It checks for a valid "Authorization" header in the Bearer format. If the token is valid, the middleware
// extracts the user ID from the token's claims and stores it in the request context.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Retrieve the token from the "Authorization" header.
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}
		// Verify that the token follows the "Bearer " format.
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		jwtToken := tokenString[7:]

		// Parse and validate the token using the provided secret.
		token, err := jwt.ParseWithClaims(jwtToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Retrieve claims from the token, specifically the subject (user ID).
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store the user ID from the claims in Gin's context and in the request context.
		c.Set(string(constants.CtxFieldUserID), claims.Subject)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.CtxFieldUserID, claims.Subject))

		// Continue to the next handler.
		c.Next()
	}
}
