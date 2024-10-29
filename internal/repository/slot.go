package repository

import (
	"context"
	"github.com/public-forge/go-gorm-unit-of-work/postgres"
	"github.com/vadymlab/slot-game/internal/interfaces"
	"github.com/vadymlab/slot-game/internal/models"
)

// slotRepository implements the ISlotRepository interface for managing
// slot game operations within the database.
type slotRepository struct{}

// AddSpin records a new spin entry in the database.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - spin: A pointer to the Spin model instance representing the spin to be recorded.
//
// Returns:
//   - An error if the transaction or spin creation fails; otherwise, nil.
func (s slotRepository) AddSpin(ctx context.Context, spin *models.Spin) error {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return err
	}

	result := tr.Provider().Create(&spin)
	if err := result.Error; err != nil {
		_ = tr.Rollback()
		return err
	}
	return tr.Commit(id)
}

// GetSpins retrieves the spin history for a specified user.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and cancellation signals.
//   - userId: The unique numeric ID of the user whose spin history is being retrieved.
//
// Returns:
//   - A slice of pointers to Spin model instances representing the user's spin history.
//   - An error if the transaction or retrieval fails; otherwise, nil.
func (s slotRepository) GetSpins(ctx context.Context, userId uint) ([]*models.Spin, error) {
	tr, _ := postgres.GetTransactionContext(ctx)
	id, err := tr.Begin()
	if err != nil {
		return nil, err
	}

	var spins []*models.Spin
	result := tr.Provider().Model(&models.Spin{}).Where("user_id = ?", userId).Find(&spins)
	if err := result.Error; err != nil {
		_ = tr.Rollback()
		return nil, err
	}
	return spins, tr.Commit(id)
}

// NewSlotRepository initializes and returns a new instance of slotRepository,
// implementing the ISlotRepository interface for slot game database operations.
func NewSlotRepository() interfaces.ISlotRepository {
	return &slotRepository{}
}
