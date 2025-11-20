package model

import (
	"time"

	"gorm.io/gorm"
)

// Email represents an email message stored in the database.
type Email struct {
	gorm.Model
	MessageID string    `gorm:"uniqueIndex;not null"`
	Subject   string
	Sender    string
	Date      time.Time
	Snippet   string
	BodyText  string    `gorm:"type:text"` // Plain text content
	BodyHTML  string    `gorm:"type:text"` // HTML content
	Summary   string    `gorm:"type:text"` // AI Generated Summary
	Sentiment string    `gorm:"size:50"`   // Positive, Neutral, Negative
	Urgency   string    `gorm:"size:50"`   // High, Medium, Low
}
