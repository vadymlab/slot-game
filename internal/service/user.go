package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	log "github.com/public-forge/go-logger"
	serviceError "github.com/vadymlab/slot-game/internal/error"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// userService implements IUserService, providing business logic for user-related actions
// such as authentication, registration, and balance management.
type userService struct {
	userRepository interfaces.IUserRepository // Repository for managing user data
}

// GetByID retrieves a user by their numeric ID.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - id: The unique numeric ID of the user.
//
// Returns:
//   - A pointer to a User model if found, or nil if not found.
//   - An error if the retrieval fails.
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	log.FromContext(ctx).Debug("Get user by id")
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByExternalID retrieves a user by their UUID identifier.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - id: A UUID representing the user's external identifier.
//
// Returns:
//   - A pointer to a User model if found, or nil if not found.
//   - An error if the retrieval fails.
func (s *userService) GetByExternalID(ctx context.Context, id *uuid.UUID) (*models.User, error) {
	log.FromContext(ctx).Debug("Get user by external id")
	user, err := s.userRepository.GetByExternalID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.FromContext(ctx).Error("User not found")
		return nil, serviceError.ErrUserNotFound
	}
	return user, nil
}

// Login authenticates a user by verifying the provided login and password.
// Logs the operation and returns an error if the user does not exist or the password is incorrect.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - login: The user's login identifier.
//   - password: The user's password for authentication.
//
// Returns:
//   - A pointer to the authenticated User model if login is successful.
//   - An error if authentication fails or the user is not found.
func (s *userService) Login(ctx context.Context, login, password string) (*models.User, error) {
	log.FromContext(ctx).Debug("Login")
	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.FromContext(ctx).Error("User not found")
		return nil, serviceError.ErrUserNotFound
	}
	if !(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil) {
		log.FromContext(ctx).Error("Password is incorrect")
		return nil, serviceError.ErrInvalidPass
	}
	return user, nil
}

// Register creates a new user with the specified login and password.
// Checks if a user with the same login already exists, hashes the password, and logs the operation.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - login: The login identifier for the new user.
//   - password: The user's password, which will be hashed before storage.
//
// Returns:
//   - A pointer to the created User model if registration is successful.
//   - An error if the registration fails or the user already exists.
func (s *userService) Register(ctx context.Context, login, password string) (*models.User, error) {
	log.FromContext(ctx).Debug("Register")
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	existUser, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	if existUser != nil {
		_ = tr.Rollback()
		return nil, serviceError.ErrUserExists
	}

	pass, err := getHash(password)
	if err != nil {
		log.FromContext(ctx).Error(err)
		_ = tr.Rollback()
		return nil, err
	}
	user := &models.User{
		Login:    login,
		Password: pass,
	}
	u, err := s.userRepository.Create(ctx, user)
	if err != nil {
		_ = tr.Rollback()
		log.FromContext(ctx).Error(err)
		return nil, err
	}
	return u, tr.Commit(id)
}

// Deposit increases a user's balance by the specified amount.
// Verifies the amount is positive, logs the operation, and performs the deposit transaction.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: The UUID representing the user's external identifier.
//   - amount: The amount to be deposited to the user's balance.
//
// Returns:
//   - A pointer to the updated balance as a float64.
//   - An error if the deposit fails or the amount is invalid.
func (s *userService) Deposit(ctx context.Context, userID *uuid.UUID, amount float64) (*float64, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}
	if amount <= 0 {
		_ = tr.Rollback()
		return nil, serviceError.ErrInvalidAmount
	}
	user, err := s.userRepository.GetByExternalID(ctx, userID)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}

	balance, err := s.userRepository.Deposit(ctx, user.ID, amount)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	return balance, tr.Commit(id)
}

// Withdraw decreases a user's balance by the specified amount.
// Checks if the user has sufficient funds, logs the operation, and performs the withdrawal transaction.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: The UUID representing the user's external identifier.
//   - amount: The amount to be withdrawn from the user's balance.
//
// Returns:
//   - A pointer to the updated balance as a float64.
//   - An error if the withdrawal fails or there are insufficient funds.
func (s *userService) Withdraw(ctx context.Context, userID *uuid.UUID, amount float64) (*float64, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetByExternalID(ctx, userID)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	if user.Balance < amount {
		_ = tr.Rollback()
		return nil, serviceError.ErrInsufficientFunds
	}
	wallet, err := s.userRepository.Withdraw(ctx, user.ID, amount)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	return wallet, tr.Commit(id)
}

// NewUserService creates and returns a new instance of userService with the given user repository.
//
// Parameters:
//   - userRepository: An implementation of IUserRepository for managing user data.
//
// Returns:
//   - A new instance of userService implementing IUserService.
func NewUserService(userRepository interfaces.IUserRepository) interfaces.IUserService {
	return &userService{
		userRepository: userRepository,
	}
}

// getHash generates a bcrypt hash from the given password string.
//
// Parameters:
//   - password: The plain text password to be hashed.
//
// Returns:
//   - The hashed password as a string.
//   - An error if hashing fails.
func getHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
