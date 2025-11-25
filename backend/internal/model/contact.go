package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contact struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID         *uuid.UUID `gorm:"type:uuid;index"` // Nullable if owned by team or org
	OrganizationID *uuid.UUID `gorm:"type:uuid;index"` // If set, available to whole org
	TeamID         *uuid.UUID `gorm:"type:uuid;index"` // If set, available to team
	IsPrivate      bool       `gorm:"default:true"`    // true: user-owned, false: shared

	Email            string `gorm:"not null"`
	Name             string
	InteractionCount int `gorm:"default:0;not null"`
	LastInteractedAt time.Time
	AvgSentiment     float64 `gorm:"type:numeric(3,2);default:0.0;not null"` // Range: -1.0 to 1.0
}
