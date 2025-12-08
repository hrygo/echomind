package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Email represents an email message stored in the database.
type Email struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	AccountID    uuid.UUID `gorm:"type:uuid;index"` // Link to EmailAccount
	MessageID    string    `gorm:"uniqueIndex;not null"`
	Subject      string
	Sender       string
	To           datatypes.JSON `gorm:"type:jsonb"` // []string
	Cc           datatypes.JSON `gorm:"type:jsonb"` // []string
	Date         time.Time
	Snippet      string
	BodyText     string         `gorm:"type:text"` // Plain text content
	BodyHTML     string         `gorm:"type:text"` // HTML content
	IsRead       bool           `gorm:"default:false"`
	Folder       string         `gorm:"size:100;default:'INBOX'"`
	Summary      string         `gorm:"type:text"`  // AI Generated Summary
	Category     string         `gorm:"size:50"`    // Work, Newsletter, Personal, etc.
	Sentiment    string         `gorm:"size:50"`    // Positive, Neutral, Negative
	Urgency      string         `gorm:"size:50"`    // High, Medium, Low
	SnoozedUntil *time.Time     `gorm:"index"`      // If set, hide from inbox until this time
	ActionItems  datatypes.JSON `gorm:"type:jsonb"` // Extracted tasks
	SmartActions datatypes.JSON `gorm:"type:jsonb"` // Structured smart actions
}
