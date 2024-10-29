package request

type BaseAuthRequest struct {
	Login    string `json:"login" validate:"required,email"`    // User's login email, required and must be a valid email format
	Password string `json:"password" validate:"required,min=8"` // User's password, required and must be at least 8 characters
}

// LoginRequest represents the request body for a user login operation.
// It includes fields for user credentials and applies validation constraints.
type LoginRequest struct {
	BaseAuthRequest
}

// RegisterRequest represents the request body for a user registration operation.
// It includes fields for user credentials and applies validation constraints.
type RegisterRequest struct {
	BaseAuthRequest
}
