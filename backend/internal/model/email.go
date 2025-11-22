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

	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	MessageID   string    `gorm:"uniqueIndex;not null"`
	Subject     string
	Sender      string
	Date        time.Time
	Snippet     string
	BodyText    string         `gorm:"type:text"`  // Plain text content
	BodyHTML    string         `gorm:"type:text"`  // HTML content
	Summary     string         `gorm:"type:text"`  // AI Generated Summary
	Category    string         `gorm:"size:50"`    // Work, Newsletter, Personal, etc.
	Sentiment   string         `gorm:"size:50"`    // Positive, Neutral, Negative
	Urgency      string         `gorm:"size:50"`    // High, Medium, Low
	ActionItems  datatypes.JSON `gorm:"type:jsonb"` // Extracted tasks
	SmartActions datatypes.JSON `gorm:"type:jsonb"` // Structured smart actions
}
