package request

// BaseAuthRequest represents the base structure for an authentication request,
// containing user credentials with validation rules for security and integrity.
type BaseAuthRequest struct {
	// Login is the user's login email address. This field is required and must
	// conform to a valid email format to ensure proper identification.
	Login string `json:"login" validate:"required,email"`

	// Password is the user's login password. This field is required and must be
	// at least 8 characters long, providing basic security against weak passwords.
	Password string `json:"password" validate:"required,min=8"`
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
