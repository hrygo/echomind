package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"not null"`
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID"`
	Members      []TeamMember `gorm:"foreignKey:TeamID"`
}
