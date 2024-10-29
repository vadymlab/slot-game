package interfaces

import (
	"context"
	"github.com/google/uuid"
	"github.com/vadymlab/slot-game/internal/models"
)

// IUserService defines service-level methods for handling user-related actions,
// including authentication, registration, and balance management operations.
type IUserService interface {
	// Login authenticates a user based on the provided login and password.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - login: The user's login identifier.
	//   - password: The user's password.
	//
	// Returns:
	//   - A pointer to a User model if authentication is successful.
	//   - An error if authentication fails or an issue occurs.
	Login(ctx context.Context, login, password string) (*models.User, error)

	// Register creates a new user account with the specified login and password.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - login: The user's desired login identifier.
	//   - password: The user's chosen password.
	//
	// Returns:
	//   - A pointer to the created User model.
	//   - An error if registration fails or an issue occurs.
	Register(ctx context.Context, login, password string) (*models.User, error)

	// GetByExternalID retrieves a user by their UUID identifier.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - id: A UUID representing the user's external identifier.
	//
	// Returns:
	//   - A pointer to a User model if found.
	//   - An error if the user is not found or if any issues occur.
	GetByExternalID(ctx context.Context, id *uuid.UUID) (*models.User, error)

	// GetByID retrieves a user by their numeric ID.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - id: The unique numeric ID of the user to retrieve.
	//
	// Returns:
	//   - A pointer to a User model if found.
	//   - An error if the user is not found or if any issues occur.
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// Deposit adds a specified amount to the balance of a user identified by their UUID.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userID: A UUID representing the user's external identifier.
	//   - amount: The amount to be deposited to the user's balance.
	//
	// Returns:
	//   - A pointer to the updated balance as a float64.
	//   - An error if the deposit fails or any issues occur.
	Deposit(ctx context.Context, userID *uuid.UUID, amount float64) (*float64, error)

	// Withdraw deducts a specified amount from the balance of a user identified by their UUID.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userId: A UUID representing the user's external identifier.
	//   - amount: The amount to be withdrawn from the user's balance.
	//
	// Returns:
	//   - A pointer to the updated balance as a float64.
	//   - An error if the withdrawal fails or any issues occur.
	Withdraw(ctx context.Context, userID *uuid.UUID, amount float64) (*float64, error)
}

// ISlotService defines service-level methods for handling slot game actions,
// including spinning and retrieving a user's spin history.
type ISlotService interface {
	RetrySpin(ctx context.Context, userID *uuid.UUID, betAmount float64) (*models.Spin, error)

	// History retrieves the spin history for a specified user.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userID: A UUID representing the user's external identifier.
	//
	// Returns:
	//   - A slice of pointers to spin models representing the user's spin history.
	//   - An error if retrieval fails or any issues occur.
	History(ctx context.Context, userID *uuid.UUID) ([]*models.Spin, error)
}
