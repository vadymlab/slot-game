package request

// BaseWalletRequest represents a base request structure for wallet transactions.
// It includes the amount to be deposited or withdrawn, with a validation constraint.
type BaseWalletRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"` // Transaction amount, required field
}

// DepositRequest represents a request to deposit funds into the user's wallet.
// It embeds BaseWalletRequest to include the amount field.
type DepositRequest struct {
	BaseWalletRequest
}

// WithdrawRequest represents a request to withdraw funds from the user's wallet.
// It embeds BaseWalletRequest to include the amount field.
type WithdrawRequest struct {
	BaseWalletRequest
}
