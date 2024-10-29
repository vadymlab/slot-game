package request

// SpinRequest represents the data required to initiate a spin in the slot game.
// The BetAmount specifies the amount of the bet placed for the spin.
type SpinRequest struct {
	BetAmount float64 `json:"bet_amount" validate:"required,gt=0"` // Bet amount, required and must be greater than 0
}
