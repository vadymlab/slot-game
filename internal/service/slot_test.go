package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	log "github.com/public-forge/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/vadymlab/slot-game/internal/config"
	error2 "github.com/vadymlab/slot-game/internal/error"
	"github.com/vadymlab/slot-game/internal/interfaces/mocks"
	"github.com/vadymlab/slot-game/internal/models"
	"testing"
)

func TestRetrySpin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	mockSlotRepo := mocks.NewMockISlotRepository(ctrl)
	mockTransactionContext := postgres.NewMockITransactionContext(ctrl)

	mockTransactionContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockTransactionContext.EXPECT().Commit(gomock.Any()).Return(nil)
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTransactionContext)

	var tr = ctx.Value(postgres.TransactionContextKey).(postgres.ITransactionContext)

	log.FromContext(ctx).Infof("Transaction context: %v", tr)

	slotConfig := &config.SlotConfig{ThreeMatchProbability: 1, TwoMatchProbability: 1, MultiplierThree: 10, MultiplierTwo: 2}

	s := NewSlotService(slotConfig, mockUserService, mockSlotRepo)

	userID := uuid.New()
	betAmount := 10.0

	mockUserService.EXPECT().GetByExternalId(ctx, &userID).Return(&models.User{
		Model: gorm.Model{
			ID: 1,
		}, Balance: 100,
	}, nil)
	mockUserService.EXPECT().Withdraw(gomock.Any(), &userID, gomock.Any()).Return(nil, nil)
	mockUserService.EXPECT().Deposit(gomock.Any(), &userID, gomock.Any()).Return(nil, nil)
	mockSlotRepo.EXPECT().AddSpin(gomock.Any(), gomock.Any()).Return(nil)

	spin, err := s.RetrySpin(ctx, &userID, betAmount)
	assert.NoError(t, err)
	assert.NotNil(t, spin)
}

func TestRetrySpin_VariousConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock services
	mockUserService := mocks.NewMockIUserService(ctrl)
	mockSlotRepo := mocks.NewMockISlotRepository(ctrl)
	mockTransactionContext := postgres.NewMockITransactionContext(ctrl)

	// Test data
	userID := uuid.New()
	betAmount := 10.0

	// Define test cases with various configurations
	testCases := []struct {
		name                  string
		threeMatchProbability float64
		twoMatchProbability   float64
		multiplierThree       float64
		multiplierTwo         float64
		expectedWin           float64
	}{
		{"ThreeMatchOnly", 1, 0, 10, 5, betAmount * 10}, // Only three-match should win
		{"TwoMatchOnly", 0, 1, 10, 5, betAmount * 5},    // Only two-match should win
		{"NoMatch", 0, 0, 10, 5, 0},                     // No matches should result in loss
		{"BothMatch", 1, 1, 10, 5, betAmount * 10},      // Three-match takes priority if both probabilities are 1
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations for transaction lifecycle
			mockTransactionContext.EXPECT().Begin().Return(uuid.New(), nil).Times(1)
			mockTransactionContext.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)
			mockTransactionContext.EXPECT().Rollback().Return(nil).Times(0)

			// Add the mocked transaction context to the context
			ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTransactionContext)

			// Set up slot configuration
			slotConfig := &config.SlotConfig{
				ThreeMatchProbability: tc.threeMatchProbability,
				TwoMatchProbability:   tc.twoMatchProbability,
				MultiplierThree:       tc.multiplierThree,
				MultiplierTwo:         tc.multiplierTwo,
			}

			// Initialize slot service
			s := NewSlotService(slotConfig, mockUserService, mockSlotRepo)

			// Expectations for user service and slot repository
			mockUserService.EXPECT().GetByExternalId(ctx, &userID).Return(&models.User{
				Model: gorm.Model{ID: 1}, Balance: 100,
			}, nil).Times(1)
			mockUserService.EXPECT().Withdraw(ctx, &userID, betAmount).Return(nil, nil).Times(1)

			// If expected win amount is greater than zero, expect a deposit
			if tc.expectedWin > 0 {
				mockUserService.EXPECT().Deposit(ctx, &userID, tc.expectedWin).Return(nil, nil).Times(1)
			}
			mockSlotRepo.EXPECT().AddSpin(ctx, gomock.Any()).Return(nil).Times(1)

			// Execute RetrySpin
			spin, err := s.RetrySpin(ctx, &userID, betAmount)

			// Assertions
			assert.NoError(t, err)
			assert.NotNil(t, spin)
			assert.Equal(t, tc.expectedWin, spin.WinAmount)
		})
	}
}

func TestRetrySpin_TemporaryError_RetrySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	mockSlotRepo := mocks.NewMockISlotRepository(ctrl)
	mockTransactionContext := postgres.NewMockITransactionContext(ctrl)

	// Set up transaction expectations
	mockTransactionContext.EXPECT().Begin().AnyTimes().Return(uuid.New(), nil)
	mockTransactionContext.EXPECT().Commit(gomock.Any()).AnyTimes().Return(nil)
	mockTransactionContext.EXPECT().Rollback().AnyTimes().Return(nil)
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTransactionContext)

	// Slot service configuration with 100% probabilities
	slotConfig := &config.SlotConfig{
		ThreeMatchProbability: 1,
		TwoMatchProbability:   1,
		MultiplierThree:       10,
		MultiplierTwo:         2,
	}
	s := NewSlotService(slotConfig, mockUserService, mockSlotRepo)

	userID := uuid.New()
	betAmount := 10.0

	// Expectations for the user and repository services
	mockUserService.EXPECT().GetByExternalId(ctx, &userID).Return(&models.User{
		Model: gorm.Model{ID: 1}, Balance: 100,
	}, nil).Times(3) // Expecting this call three times due to retries
	mockUserService.EXPECT().Withdraw(ctx, &userID, betAmount).Return(nil, error2.ErrInsufficientFunds).Times(2)
	mockUserService.EXPECT().Withdraw(ctx, &userID, betAmount).Return(nil, nil).Times(1)
	mockUserService.EXPECT().Deposit(ctx, &userID, gomock.Any()).Return(nil, nil).Times(1)
	mockSlotRepo.EXPECT().AddSpin(ctx, gomock.Any()).Return(nil).Times(1)

	// Execute RetrySpin
	spin, err := s.RetrySpin(ctx, &userID, betAmount)

	// Assertions to verify retry behavior and results
	assert.NoError(t, err)
	assert.NotNil(t, spin)
	assert.Equal(t, betAmount*slotConfig.MultiplierThree, spin.WinAmount)
}

func TestHistory_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	mockUserService := mocks.NewMockIUserService(ctrl)
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)

	userID := uuid.New()
	expectedErr := errors.New("user not found")

	// Expectations
	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserService.EXPECT().GetByExternalId(ctx, &userID).Return(nil, expectedErr)
	mockTxContext.EXPECT().Rollback().Return(nil)

	// Instantiate the service
	service := NewSlotService(nil, mockUserService, nil)

	// Act
	history, err := service.History(ctx, &userID)

	// Assert
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, history)
}

func TestHistory_GetSpinsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	mockUserService := mocks.NewMockIUserService(ctrl)
	mockSlotRepo := mocks.NewMockISlotRepository(ctrl)
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)

	userID := uuid.New()
	mockUser := &models.User{Model: gorm.Model{ID: 1}}
	expectedErr := errors.New("failed to fetch spins")

	// Expectations
	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserService.EXPECT().GetByExternalId(ctx, &userID).Return(mockUser, nil)
	mockSlotRepo.EXPECT().GetSpins(ctx, mockUser.ID).Return(nil, expectedErr)
	mockTxContext.EXPECT().Rollback().Return(nil)

	// Instantiate the service
	service := NewSlotService(nil, mockUserService, mockSlotRepo)

	// Act
	history, err := service.History(ctx, &userID)

	// Assert
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, history)
}
func TestHistory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	mockUserService := mocks.NewMockIUserService(ctrl)
	mockSlotRepo := mocks.NewMockISlotRepository(ctrl)
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)

	// Test Data
	userID := uuid.New()
	mockUser := &models.User{Model: gorm.Model{ID: 1}}
	mockHistory := []*models.Spin{{UserID: mockUser.ID}}

	// Expectations
	mockTxContext.EXPECT().Begin().Return(uuid.New(), nil)
	mockUserService.EXPECT().GetByExternalId(ctx, &userID).Return(mockUser, nil)
	mockSlotRepo.EXPECT().GetSpins(ctx, mockUser.ID).Return(mockHistory, nil)
	mockTxContext.EXPECT().Commit(gomock.Any()).Return(nil)

	// Instantiate the service
	service := NewSlotService(nil, mockUserService, mockSlotRepo)

	// Act
	history, err := service.History(ctx, &userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockHistory, history)
}

func TestHistory_BeginTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockTxContext := postgres.NewMockITransactionContext(ctrl)
	ctx := context.WithValue(context.Background(), postgres.TransactionContextKey, mockTxContext)
	expectedErr := errors.New("failed to begin transaction")

	// Expectations
	mockTxContext.EXPECT().Begin().Return(uuid.Nil, expectedErr)

	// Instantiate the service
	service := NewSlotService(nil, nil, nil)

	uid := uuid.New()
	// Act
	history, err := service.History(ctx, &uid)

	// Assert
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, history)
}
