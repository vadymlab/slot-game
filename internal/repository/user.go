package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/models"
)

// userRepository implements IUserRepository interface for accessing
// and managing user-related data in the database.
type userRepository struct{}

// GetById retrieves a user by their numeric ID.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - uid: The unique numeric ID of the user.
//
// Returns:
//   - A pointer to a User model if found, or nil if not found.
//   - An error if the retrieval fails.
func (r *userRepository) GetById(ctx context.Context, uid uint) (*models.User, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	result := tr.Provider().Model(&models.User{}).Where("id = ?", uid).First(user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		_ = tr.Rollback()
		return nil, err
	}
	return user, tr.Commit(id)
}

// GetByExternalId retrieves a user by their UUID identifier.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: A UUID representing the user's external identifier.
//
// Returns:
//   - A pointer to a User model if found, or nil if not found.
//   - An error if the retrieval fails.
func (r *userRepository) GetByExternalId(ctx context.Context, userId *uuid.UUID) (*models.User, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	result := tr.Provider().Model(&models.User{}).Where("external_id = ?", userId).First(user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		_ = tr.Rollback()
		return nil, err
	}
	return user, tr.Commit(id)
}

// Create inserts a new user record into the database.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - user: A pointer to a User model instance representing the new user.
//
// Returns:
//   - A pointer to the created User model.
//   - An error if the creation fails.
func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	result := tr.Provider().Create(&user)
	if err := result.Error; err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	return user, tr.Commit(id)
}

// GetByLogin retrieves a user by their login name.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - login: The login name of the user.
//
// Returns:
//   - A pointer to a User model if found, or nil if not found.
//   - An error if the retrieval fails.
func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	result := tr.Provider().Model(&models.User{}).Where("login = ?", login).First(user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		_ = tr.Rollback()
		return nil, err
	}
	return user, tr.Commit(id)
}

// Deposit increases the balance of a specified user.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: The unique numeric ID of the user.
//   - amount: The amount to be added to the user's balance.
//
// Returns:
//   - A pointer to the updated balance as a float64.
//   - An error if the update fails.
func (r *userRepository) Deposit(ctx context.Context, userId uint, amount float64) (*float64, error) {
	return r.updateBalance(ctx, userId, amount)
}

// Withdraw decreases the balance of a specified user.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: The unique numeric ID of the user.
//   - amount: The amount to be deducted from the user's balance.
//
// Returns:
//   - A pointer to the updated balance as a float64.
//   - An error if the update fails.
func (r *userRepository) Withdraw(ctx context.Context, userId uint, amount float64) (*float64, error) {
	return r.updateBalance(ctx, userId, -amount)
}

// updateBalance modifies the balance of a specified user by the given amount.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: The unique numeric ID of the user.
//   - amount: The amount to be added to or deducted from the balance.
//
// Returns:
//   - A pointer to the updated balance as a float64.
//   - An error if the update fails.
func (r *userRepository) updateBalance(ctx context.Context, userId uint, amount float64) (*float64, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	result := tr.Provider().Model(user).Where("id = ?", userId).First(user)
	if err := result.Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			_ = tr.Rollback()
			return nil, err
		}
	}
	user.Balance += amount

	result = tr.Provider().Model(&user).
		Where("id = ?", userId).
		Update("balance", user.Balance).
		Select("balance").
		Scan(&user)
	if err := result.Error; err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	return &user.Balance, tr.Commit(id)
}

// NewUserRepository creates and returns a new instance of userRepository.
func NewUserRepository() interfaces.IUserRepository {
	return &userRepository{}
}
