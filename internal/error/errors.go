package error

// Predefined user-related errors.
var (
	ErrUserNotFound      = &UserNotFound{}      // Error for when a user cannot be found
	ErrUserExists        = &UserAlreadyExists{} // Error for when a user already exists during registration
	ErrInvalidPass       = &InvalidPassword{}   // Error for when user credentials are incorrect
	ErrInsufficientFunds = &InefficientFunds{}  // Error for when a user has insufficient funds for a transaction
	ErrInvalidAmount     = &InefficientFunds{}  // Error for when a transaction amount is invalid
)

// UserNotFound represents an error for when a requested user does not exist.
type UserNotFound struct{}

// UserAlreadyExists represents an error for when a user with the specified login already exists.
type UserAlreadyExists struct{}

// InvalidPassword represents an error for incorrect user credentials.
type InvalidPassword struct{}

// InefficientFunds represents an error for insufficient funds during a transaction.
type InefficientFunds struct{}

// InvalidAmount represents an error for an invalid transaction amount.
type InvalidAmount struct{}

// Error returns the error message for UserNotFound.
func (cs UserNotFound) Error() string {
	return "user not found"
}

// Error returns the error message for UserAlreadyExists.
func (cs UserAlreadyExists) Error() string {
	return "user with this login already exists"
}

// Error returns the error message for InvalidPassword.
func (cs InvalidPassword) Error() string {
	return "wrong credentials"
}

// Error returns the error message for InefficientFunds.
func (cs InefficientFunds) Error() string {
	return "insufficient funds"
}

// Error returns the error message for InvalidAmount.
func (cs InvalidAmount) Error() string {
	return "invalid amount"
}
