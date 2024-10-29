package response

import "github.com/vadymlab/slot-game/internal/models"

// SpinResponse represents the response returned after a spin is completed,
// containing information about the amount won in that spin.
type SpinResponse struct {
	WinAmount float64 `json:"win_amount"` // The amount the user won on this spin
}

// SpinHistoryResponse represents a structured response for a user's spin history.
// It includes essential details such as the bet amount, win amount, and the date of each spin.
type SpinHistoryResponse struct {
	BetAmount float64 `json:"bet_amount"` // The amount the user bet on this spin
	WinAmount float64 `json:"win_amount"` // The amount the user won on this spin
	Date      string  `json:"date"`       // The date and time of this spin, formatted as "YYYY-MM-DD HH:MM:SS"
}

// SpinFromModel creates a SpinResponse instance from a Spin model.
// This function is used to generate a response object with the winning amount from a spin.
//
// Parameters:
//   - model: A pointer to a models.Spin instance containing the spin data.
//
// Returns:
//
//	A pointer to a SpinResponse instance with the win amount mapped from the input model.
func SpinFromModel(model *models.Spin) *SpinResponse {
	return &SpinResponse{
		WinAmount: model.WinAmount,
	}
}

// SpinHistoryFromModel converts a Spin model instance to a SpinHistoryResponse instance.
// This function is used to create a serializable response object for a single spin record in history.
//
// Parameters:
//   - model: A pointer to a models.Spin instance containing the original spin data.
//
// Returns:
//
//	A pointer to a SpinHistoryResponse instance containing the mapped data from the input model.
func SpinHistoryFromModel(model *models.Spin) *SpinHistoryResponse {
	return &SpinHistoryResponse{
		BetAmount: model.BetAmount,
		WinAmount: model.WinAmount,
		Date:      model.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// SpinHistoryFromModels converts a slice of Spin model instances to a slice of SpinHistoryResponse instances.
// This function generates a collection of serializable response objects for multiple spin records in history.
//
// Parameters:
//   - models: A slice of pointers to models.Spin instances, each representing a spin record.
//
// Returns:
//
//	A slice of pointers to SpinHistoryResponse instances, with each entry containing mapped data
//	from a corresponding input model.
func SpinHistoryFromModels(models []*models.Spin) []*SpinHistoryResponse {
	var res []*SpinHistoryResponse
	for _, model := range models {
		res = append(res, SpinHistoryFromModel(model))
	}
	return res
}
