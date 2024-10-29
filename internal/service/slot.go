package service

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	error2 "github.com/vadymlab/slot-game/internal/error"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	log "github.com/public-forge/go-logger"
	"github.com/vadymlab/slot-game/internal/config"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/models"
)

// symbols defines the available slot machine symbols.
var symbols = []string{"A", "B", "C", "D"}

// slotService implements ISlotService, providing slot game logic and methods.
type slotService struct {
	config         *config.SlotConfig         // Slot configuration settings
	userService    interfaces.IUserService    // Service for managing user-related operations
	slotRepository interfaces.ISlotRepository // Repository for managing slot spin records
	rng            *rand.Rand                 // Custom random number generator for reproducibility
	backoff        *backoff.ExponentialBackOff
}

// History retrieves the spin history for a specified user.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: A UUID representing the user's external identifier.
//
// Returns:
//   - A slice of pointers to Spin models representing the user's spin history.
//   - An error if the transaction or retrieval fails; otherwise, nil.
func (s *slotService) History(ctx context.Context, userID *uuid.UUID) ([]*models.Spin, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}
	user, err := s.userService.GetByExternalID(ctx, userID)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	history, err := s.slotRepository.GetSpins(ctx, user.ID)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	return history, tr.Commit(id)
}

// RetrySpin performs a slot spin operation for a user with a retry mechanism.
// The function attempts to execute a spin with a specified bet amount, automatically
// retrying on errors except when the error is due to insufficient funds.
//
// Parameters:
//   - ctx: A context.Context for request-scoped values and cancelation signals.
//   - userId: A UUID pointer representing the unique identifier of the user.
//   - betAmount: A float64 representing the bet amount for the spin.
//
// Returns:
//   - *models.Spin: A pointer to a Spin object containing the spin details if successful.
//   - error: An error indicating failure reason, or nil if the spin succeeds.
//
// Workflow:
//  1. Defines the `operation` function, which performs the spin and retries
//     unless the error is due to insufficient funds, in which case it immediately returns.
//  2. The `backoff.Retry` function is called, which retries `operation` based on
//     the backoff configuration specified in `s.backoff`.
//  3. Logs warning messages for insufficient funds and error messages for retries
//     that exceed the allowed backoff configuration.
//
// Logging:
//
//	Logs messages at the appropriate levels (warning, error, and debug) based on
//	the retry status and outcomes, including retry counts and error details.
//
// Example usage:
//
//	spin, err := slotService.RetrySpin(ctx, &userId, betAmount)
//	if err != nil {
//	    // Handle error
//	}
//	// Process spin result
func (s *slotService) RetrySpin(ctx context.Context, userID *uuid.UUID, betAmount float64) (*models.Spin, error) {
	var spin *models.Spin
	operation := func() error {
		var err error
		spin, err = s.spin(ctx, userID, betAmount)
		if err != nil {
			if errors.Is(err, error2.ErrInsufficientFunds) {
				log.FromContext(ctx).Warnf("RetrySpin encountered error: %v", err)
				return err
			}
			return backoff.Permanent(err)
		}
		return nil
	}

	// Run the operation with retries
	err := backoff.Retry(operation, s.backoff)
	if err != nil {
		log.FromContext(ctx).Errorf("RetrySpin failed after %v retries: %v", s.backoff.MaxElapsedTime, err)
		return nil, err
	}

	log.FromContext(ctx).Debugf("RetrySpin succeeded after %v retries", s.backoff.GetElapsedTime())
	return spin, nil
}

// spin initiates a spin for the slot machine with a specified bet amount,
// calculates the payout, and updates the user's balance.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: A UUID representing the user's external identifier.
//   - betAmount: The amount of the bet placed for the spin.
//
// Returns:
//   - A pointer to a spin model representing the spin result.
//   - An error if the spin process or transaction fails; otherwise, nil.
func (s *slotService) spin(ctx context.Context, userID *uuid.UUID, betAmount float64) (*models.Spin, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}
	user, err := s.userService.GetByExternalID(ctx, userID)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	_, err = s.userService.Withdraw(ctx, userID, betAmount)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}

	payout := s.calculatePayout(betAmount)
	if payout > 0 {
		_, err = s.userService.Deposit(ctx, userID, payout)
		if err != nil {
			_ = tr.Rollback()
			return nil, err
		}
	}

	spin := &models.Spin{
		UserID:    user.ID,
		BetAmount: betAmount,
		WinAmount: payout,
	}
	err = s.slotRepository.AddSpin(ctx, spin)
	if err != nil {
		_ = tr.Rollback()
		return nil, err
	}

	log.FromContext(ctx).Infof("spin result: %+v", spin)
	return spin, tr.Commit(id)
}

// calculatePayout determines the payout based on the bet amount and spin result.
// It applies predefined multipliers and winning probabilities for symbol matches.
//
// Parameters:
//   - betAmount: The amount of the bet placed for the spin.
//
// Returns:
//   - The calculated payout amount, based on the match conditions and probabilities.
func (s *slotService) calculatePayout(betAmount float64) float64 {
	// Generate random symbols for the spin result.
	spinResult := []string{
		symbols[s.rng.Intn(len(symbols))],
		symbols[s.rng.Intn(len(symbols))],
		symbols[s.rng.Intn(len(symbols))],
	}

	// Check for a three-symbol match based on ThreeMatchProbability.
	// If probability conditions are met, create a matching three-symbol result
	// and return the payout calculated with MultiplierThree.
	if s.rng.Float64() <= s.config.ThreeMatchProbability {
		spinResult[1] = spinResult[0]
		spinResult[2] = spinResult[0]
		return betAmount * s.config.MultiplierThree
	}

	// Check for a two-symbol match based on TwoMatchProbability.
	// If probability conditions are met, create a matching two-symbol result
	// and return the payout calculated with MultiplierTwo.
	if s.rng.Float64() <= s.config.TwoMatchProbability {
		spinResult[1] = spinResult[0]
		return betAmount * s.config.MultiplierTwo
	}

	// No matching symbols result in a loss with zero payout.
	return 0
}

// NewSlotService creates and returns a new instance of slotService.
//
// Parameters:
//   - config: SlotConfig containing slot game settings.
//   - userService: UserService for managing user-related operations.
//   - slotRepository: SlotRepository for handling spin records.
//
// Returns:
//   - An instance of slotService implementing ISlotService.
func NewSlotService(
	config *config.SlotConfig,
	userService interfaces.IUserService,
	slotRepository interfaces.ISlotRepository,
) interfaces.ISlotService {
	return &slotService{
		config:         config,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		userService:    userService,
		slotRepository: slotRepository,
		backoff: backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(500*time.Millisecond),
			backoff.WithMaxElapsedTime(2*time.Second),
			backoff.WithMultiplier(1.5),
		),
	}
}
