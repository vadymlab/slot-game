package response

// DepositResponse represents the response body for a successful deposit transaction.
// It includes the updated wallet balance after the deposit.
type DepositResponse struct {
	Balance float64 `json:"balance"` // Updated wallet balance after the deposit transaction
}

// WithdrawResponse represents the response body for a successful withdrawal transaction.
// It includes the updated wallet balance after the withdrawal.
type WithdrawResponse struct {
	Balance float64 `json:"balance"` // Updated wallet balance after the withdrawal transaction
}
