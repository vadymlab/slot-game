package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// User represents a registered user in the system, storing essential
// account details such as login credentials, balance, and unique identifiers.
type User struct {
	gorm.Model
	ExternalID *uuid.UUID `gorm:"column:external_id;type:uuid;default:uuid_generate_v4();unique;not null"` // Unique UUID for external identification
	Login      string     `gorm:"column:login;unique;not null"`                                            // Unique login name for the user
	Password   string     `gorm:"column:password;not null"`                                                // User's hashed password
	Balance    float64    `gorm:"column:balance;default:null"`                                             // User's current wallet balance
}

// TableName sets the table name for the User model explicitly.
func (User) TableName() string {
	return "users"
}
