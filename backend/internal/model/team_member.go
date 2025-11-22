package model

import (
	"time"

	"github.com/google/uuid"
)

type TeamMember struct {
	TeamID    uuid.UUID        `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Role      OrganizationRole `gorm:"type:varchar(20);default:'member'"` // Reusing OrganizationRole for now
	JoinedAt  time.Time        `gorm:"autoCreateTime"`

	// Associations
	User User `gorm:"foreignKey:UserID"`
	Team Team `gorm:"foreignKey:TeamID"`
}
