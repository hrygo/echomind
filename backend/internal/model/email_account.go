package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailAccount represents a user's configured email account for IMAP synchronization.
type EmailAccount struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	
	Email             string    `gorm:"not null"`          // The email address (display/login if different from username)
	ServerAddress     string    `gorm:"not null"`          // e.g., imap.gmail.com
	ServerPort        int       `gorm:"not null;default:993"`
	Username          string    `gorm:"not null"`          // IMAP Login Username
	
	EncryptedPassword string    `gorm:"type:text;not null"` // Base64 encoded ciphertext
	
	IsConnected       bool      `gorm:"default:false"`     // Status flag: true if last connection attempt was successful
	LastSyncAt        *time.Time                           // Timestamp of last successful sync
	ErrorMessage      string    `gorm:"type:text"`         // Stores the error message from the last failed connection/sync attempt
}
