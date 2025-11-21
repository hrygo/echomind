package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"` // bcrypt hash
	Name         string    `gorm:"type:varchar(100)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
