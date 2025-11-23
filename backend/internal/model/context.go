package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Context represents a smart context (project/topic) for organizing emails and tasks.
type Context struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name         string         `gorm:"type:varchar(100);not null"`
	Color        string         `gorm:"type:varchar(20);default:'blue'"`
	Keywords     datatypes.JSON `gorm:"type:jsonb"` // []string
	Stakeholders datatypes.JSON `gorm:"type:jsonb"` // []string (email addresses)
}

// EmailContext represents the many-to-many relationship between Emails and Contexts.
type EmailContext struct {
	EmailID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	ContextID uuid.UUID `gorm:"primaryKey;type:uuid"`
}
