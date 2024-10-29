package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	"github.com/stretchr/testify/assert"
	serviceError "github.com/vadymlab/slot-game/internal/error"
	"github.com/vadymlab/slot-game/internal/interfaces/mocks"
	"github.com/vadymlab/slot-game/internal/models"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestGetById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	userID := uint(1)
	expectedUser := &models.User{Model: gorm.Model{ID: userID}}

	// Expectations
	mockUserRepo.EXPECT().GetById(ctx, userID).Return(expectedUser, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetById_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	userID := uint(1)
	expectedErr := errors.New("user not found")

	// Expectations
	mockUserRepo.EXPECT().GetById(ctx, userID).Return(nil, expectedErr)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByID(ctx, userID)

	// Assert
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, user)
}

func TestGetById_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	userID := uint(1)
	expectedErr := errors.New("repository error")

	// Expectations
	mockUserRepo.EXPECT().GetById(ctx, userID).Return(nil, expectedErr)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByID(ctx, userID)

	// Assert
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, user)
}

func TestGetById_EmptyUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	userID := uint(1)
	emptyUser := &models.User{} // Empty user struct

	// Expectations
	mockUserRepo.EXPECT().GetById(ctx, userID).Return(emptyUser, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, emptyUser, user)
	assert.Equal(t, uint(0), user.ID) // Проверка, что ID равен 0
}

func TestGetByExternalId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	externalID := uuid.New()
	expectedUser := &models.User{ExternalID: &externalID}

	// Expectations
	mockUserRepo.EXPECT().GetByExternalId(ctx, &externalID).Return(expectedUser, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByExternalID(ctx, &externalID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetByExternalId_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	externalID := uuid.New()

	// Expectations
	mockUserRepo.EXPECT().GetByExternalId(ctx, &externalID).Return(nil, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByExternalID(ctx, &externalID)

	// Assert
	assert.ErrorIs(t, err, serviceError.ErrUserNotFound)
	assert.Nil(t, user)
}

func TestGetByExternalId_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	externalID := uuid.New()
	expectedError := errors.New("repository error")

	// Expectations
	mockUserRepo.EXPECT().GetByExternalId(ctx, &externalID).Return(nil, expectedError)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.GetByExternalID(ctx, &externalID)

	// Assert
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, user)
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	login := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedUser := &models.User{Login: login, Password: string(hashedPassword)}

	// Expectations
	mockUserRepo.EXPECT().GetByLogin(ctx, login).Return(expectedUser, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.Login(ctx, login, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	login := "nonexistentuser"
	password := "password123"

	// Expectations
	mockUserRepo.EXPECT().GetByLogin(ctx, login).Return(nil, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.Login(ctx, login, password)

	// Assert
	assert.ErrorIs(t, err, serviceError.ErrUserNotFound)
	assert.Nil(t, user)
}

func TestLogin_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	login := "testuser"
	password := "password123"
	expectedError := errors.New("repository error")

	// Expectations
	mockUserRepo.EXPECT().GetByLogin(ctx, login).Return(nil, expectedError)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.Login(ctx, login, password)

	// Assert
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, user)
}

func TestLogin_IncorrectPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	ctx := context.Background()
	login := "testuser"
	password := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedUser := &models.User{Login: login, Password: string(hashedPassword)}

	// Expectations
	mockUserRepo.EXPECT().GetByLogin(ctx, login).Return(expectedUser, nil)

	// Instantiate the service
	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.Login(ctx, login, wrongPassword)

	// Assert
	assert.ErrorIs(t, err, serviceError.ErrInvalidPass)
	assert.Nil(t, user)
}
func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	// Arrange
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	login := "newuser"
	password := "password123"

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockTxContext.EXPECT().Commit(gomock.Any()).Return(nil)
	mockUserRepo.EXPECT().GetByLogin(ctx, login).Return(nil, nil)
	// Using AssignableToTypeOf to ignore the specific password hash value
	mockUserRepo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(&models.User{Login: login})).Return(&models.User{Login: login}, nil)

	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.Register(ctx, login, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, login, user.Login)
}

func TestRegister_UserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	// Arrange
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	login := "existinguser"
	password := "password123"
	existingUser := &models.User{Login: login}

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockTxContext.EXPECT().Rollback().Return(nil)
	mockUserRepo.EXPECT().GetByLogin(ctx, login).Return(existingUser, nil)

	service := NewUserService(mockUserRepo)

	// Act
	user, err := service.Register(ctx, login, password)

	// Assert
	assert.ErrorIs(t, err, serviceError.ErrUserExists)
	assert.Nil(t, user)
}

func TestDeposit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	userID := uuid.New()
	amount := 100.0
	initialBalance := 50.0
	expectedBalance := initialBalance + amount

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(&models.User{
		Model: gorm.Model{ID: 1},
	}, nil)
	mockUserRepo.EXPECT().Deposit(ctx, uint(1), amount).Return(&expectedBalance, nil)
	mockTxContext.EXPECT().Commit(gomock.Any()).Return(nil)

	service := userService{
		userRepository: mockUserRepo,
	}
	balance, err := service.Deposit(ctx, &userID, amount)

	assert.NoError(t, err)
	assert.Equal(t, &expectedBalance, balance)
}

func TestDeposit_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	userID := uuid.New()
	amount := 0.0

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockTxContext.EXPECT().Rollback().Return(nil)

	service := userService{
		userRepository: mockUserRepo,
	}
	balance, err := service.Deposit(ctx, &userID, amount)

	assert.Nil(t, balance)
	assert.ErrorIs(t, err, serviceError.ErrInvalidAmount)
}

func TestDeposit_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	userID := uuid.New()
	amount := 100.0

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(nil, serviceError.ErrUserNotFound)
	mockTxContext.EXPECT().Rollback().Return(nil)

	service := userService{
		userRepository: mockUserRepo,
	}
	balance, err := service.Deposit(ctx, &userID, amount)

	assert.Nil(t, balance)
	assert.ErrorIs(t, err, serviceError.ErrUserNotFound)
}

func TestDeposit_DepositError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	userID := uuid.New()
	amount := 100.0

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(&models.User{
		Model: gorm.Model{ID: 1},
	}, nil)
	mockUserRepo.EXPECT().Deposit(ctx, uint(1), amount).Return(nil, errors.New("deposit error"))
	mockTxContext.EXPECT().Rollback().Return(nil)

	service := userService{
		userRepository: mockUserRepo,
	}
	balance, err := service.Deposit(ctx, &userID, amount)

	assert.Nil(t, balance)
	assert.EqualError(t, err, "deposit error")
}

func TestDeposit_CommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTxContext := postgres.NewMockITransactionContext(ctrl)

	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	userID := uuid.New()
	amount := 100.0
	initialBalance := 50.0
	expectedBalance := initialBalance + amount

	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(&models.User{
		Model: gorm.Model{ID: 1},
	}, nil)
	mockUserRepo.EXPECT().Deposit(ctx, uint(1), amount).Return(&expectedBalance, nil)
	mockTxContext.EXPECT().Commit(gomock.Any()).Return(errors.New("commit error"))

	service := userService{
		userRepository: mockUserRepo,
	}
	_, err := service.Deposit(ctx, &userID, amount)

	assert.EqualError(t, err, "commit error")
}
func TestWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)

	// Set up expected behavior for the transaction context
	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil).Times(1)
	mockTxContext.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)

	// Add the mocked transaction context to the context
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)

	// Test parameters
	userID := uuid.New()
	amount := 50.0
	user := &models.User{
		Model:   gorm.Model{ID: 1},
		Balance: 100.0,
	}
	expectedBalance := user.Balance - amount // Calculate expected balance

	// Set up expectations for repository methods
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(user, nil).Times(1)
	mockUserRepo.EXPECT().Withdraw(ctx, user.ID, amount).Return(&expectedBalance, nil).Times(1)

	service := userService{
		userRepository: mockUserRepo,
	}

	// Execute Withdraw
	balance, err := service.Withdraw(ctx, &userID, amount)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, &expectedBalance, balance) // Ensure returned balance matches expectation
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)

	// Expect calls for transaction context methods
	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockTxContext.EXPECT().Rollback().Return(nil)

	// Create a context with the mocked transaction context
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)

	// Test parameters
	userID := uuid.New()
	amount := 150.0
	user := &models.User{
		Model:   gorm.Model{ID: 1},
		Balance: 100.0,
	}

	// Set up expectations for repository methods
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(user, nil)

	service := userService{
		userRepository: mockUserRepo,
	}

	// Execute the method being tested
	wallet, err := service.Withdraw(ctx, &userID, amount)

	// Verify results
	assert.Nil(t, wallet)
	assert.ErrorIs(t, err, serviceError.ErrInsufficientFunds)
}

func TestWithdraw_ErrorInWithdrawRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)

	// Expect calls for transaction context methods
	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockTxContext.EXPECT().Rollback().Return(nil)

	// Create a context with the mocked transaction context
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)

	// Test parameters
	userID := uuid.New()
	amount := 50.0
	user := &models.User{
		Model:   gorm.Model{ID: 1},
		Balance: 100.0,
	}

	// Set up expectations for repository methods
	expectedError := errors.New("repository error")
	mockUserRepo.EXPECT().GetByExternalId(ctx, &userID).Return(user, nil)
	mockUserRepo.EXPECT().Withdraw(ctx, user.ID, amount).Return(nil, expectedError)

	service := userService{
		userRepository: mockUserRepo,
	}

	// Execute the method being tested
	wallet, err := service.Withdraw(ctx, &userID, amount)

	// Verify results
	assert.Nil(t, wallet)
	assert.ErrorIs(t, err, expectedError)
}
