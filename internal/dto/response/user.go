package response

import (
	"github.com/google/uuid"
	"github.com/vadymlab/slot-game/internal/models"
)

// LoginResponse represents the response body for a successful login operation.
// This includes a JWT token, which is required for subsequent authenticated requests.
type LoginResponse struct {
	Token string `json:"token"` // JWT token for the authenticated user
}

// ProfileResponse represents the response body for retrieving a user's profile information.
// It includes the user's unique identifier, login, and wallet balance.
type ProfileResponse struct {
	ID      *uuid.UUID `json:"id"`      // Unique identifier for the user
	Login   string     `json:"login"`   // User's login name
	Balance float64    `json:"balance"` // User's current wallet balance
}

// RegisterResponse represents the response body for a successful user registration.
// It includes the user's unique identifier and login information.
type RegisterResponse struct {
	ID    *uuid.UUID `json:"id"`    // Unique identifier for the newly registered user
	Login string     `json:"login"` // Login name for the newly registered user
}

// RegisterFromModel creates a RegisterResponse instance from a User model.
// This function is used to generate a serializable response object upon successful user registration.
//
// Parameters:
//   - user: A pointer to a models.User instance containing user data.
//
// Returns:
//
//	A pointer to a RegisterResponse instance containing the user's ID and login.
func RegisterFromModel(user *models.User) *RegisterResponse {
	return &RegisterResponse{
		ID:    user.ExternalID,
		Login: user.Login,
	}
}
