package model

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationRole string

const (
	OrgRoleOwner  OrganizationRole = "owner"
	OrgRoleAdmin  OrganizationRole = "admin"
	OrgRoleMember OrganizationRole = "member"
)

type OrganizationMember struct {
	OrganizationID uuid.UUID        `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Role           OrganizationRole `gorm:"type:varchar(20);default:'member'"`
	JoinedAt       time.Time        `gorm:"autoCreateTime"`

	// Associations
	User         User         `gorm:"foreignKey:UserID"`
	Organization Organization `gorm:"foreignKey:OrganizationID"`
}
