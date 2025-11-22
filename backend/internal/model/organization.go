package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"not null"`
	Slug        string    `gorm:"uniqueIndex;not null"` // For URL routing /orgs/:slug
	OwnerID     uuid.UUID `gorm:"type:uuid;not null"`   // The super admin
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Associations
	Members []OrganizationMember `gorm:"foreignKey:OrganizationID"`
}
