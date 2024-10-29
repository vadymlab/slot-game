package interfaces

import (
	"context"
	"github.com/google/uuid"
	"github.com/vadymlab/slot-game/internal/models"
)

// IUserRepository defines methods for user data operations in the repository layer.
type IUserRepository interface {
	// GetByLogin retrieves a user by their login name.
	// Returns the user if found, otherwise returns nil and an error if any issues occur.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - login: The login name of the user to retrieve.
	//
	// Returns:
	//   - A pointer to a User model if found.
	//   - An error if any issues occur during retrieval.
	GetByLogin(ctx context.Context, login string) (*models.User, error)

	// Create adds a new user to the repository.
	// Takes a User model as input and returns the created user and an error if any issues occur.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - user: A pointer to a User model representing the new user to add.
	//
	// Returns:
	//   - A pointer to the created User model.
	//   - An error if any issues occur during creation.
	Create(ctx context.Context, user *models.User) (*models.User, error)

	// GetByExternalId retrieves a user by their UUID identifier.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - id: A UUID representing the external identifier of the user.
	//
	// Returns:
	//   - A pointer to a User model if found.
	//   - An error if any issues occur during retrieval.
	GetByExternalId(ctx context.Context, id *uuid.UUID) (*models.User, error)

	// GetById retrieves a user by their numeric ID.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - id: The unique numeric ID of the user to retrieve.
	//
	// Returns:
	//   - A pointer to a User model if found.
	//   - An error if any issues occur during retrieval.
	GetById(ctx context.Context, id uint) (*models.User, error)

	// Deposit increases the balance of a specified user by the given amount.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userId: The unique numeric ID of the user to deposit funds into.
	//   - amount: The amount to deposit to the user's balance.
	//
	// Returns:
	//   - A pointer to the updated balance as a float64.
	//   - An error if any issues occur during the deposit.
	Deposit(ctx context.Context, userId uint, amount float64) (*float64, error)

	// Withdraw decreases the balance of a specified user by the given amount.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userId: The unique numeric ID of the user to withdraw funds from.
	//   - amount: The amount to withdraw from the user's balance.
	//
	// Returns:
	//   - A pointer to the updated balance as a float64.
	//   - An error if any issues occur during the withdrawal.
	Withdraw(ctx context.Context, userId uint, amount float64) (*float64, error)
}

// IWalletRepository defines methods for wallet-related data operations in the repository layer.
type IWalletRepository interface {
	// GetBalance retrieves the balance of a specified user.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userId: The unique numeric ID of the user whose balance is being retrieved.
	//
	// Returns:
	//   - The user's balance as a float64.
	//   - An error if any issues occur during retrieval.
	GetBalance(ctx context.Context, userId uint) (float64, error)
}

// ISlotRepository defines methods for slot game data operations in the repository layer.
type ISlotRepository interface {
	// AddSpin records a new spin for a user in the repository.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - spin: A pointer to a Spin model containing spin data to be recorded.
	//
	// Returns:
	//   - An error if any issues occur during recording of the spin.
	AddSpin(ctx context.Context, spin *models.Spin) error

	// GetSpins retrieves a user's spin history from the repository.
	//
	// Parameters:
	//   - ctx: Context for managing request-scoped values and cancellation signals.
	//   - userId: The unique numeric ID of the user whose spin history is being retrieved.
	//
	// Returns:
	//   - A slice of pointers to Spin models representing the user's spin history.
	//   - An error if any issues occur during retrieval.
	GetSpins(ctx context.Context, userId uint) ([]*models.Spin, error)
}
