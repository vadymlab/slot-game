package models

import "github.com/jinzhu/gorm"

// Spin represents a spin entry linked to a user. Each spin stores the bet amount,
// win amount, and a reference to the user who initiated the spin.
type Spin struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`                                                         // Foreign key to the User model
	BetAmount float64 `gorm:"column:bet_amount;not null"`                                       // The amount bet for this spin
	WinAmount float64 `gorm:"column:win_amount;not null"`                                       // The amount won for this spin
	User      User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Association to the User, with update and delete constraints
}

// TableName sets the table name for the Spin model explicitly.
func (Spin) TableName() string {
	return "spins"
}
