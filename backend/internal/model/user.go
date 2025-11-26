package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"` // bcrypt hash
	Name         string    `gorm:"type:varchar(100)"`
	Role         string    `gorm:"type:varchar(50);default:'manager';not null"` // Add Role field
	
	// WeChat Integration
	WeChatOpenID  string `gorm:"type:varchar(64);uniqueIndex"` // Unique per app
	WeChatUnionID string `gorm:"type:varchar(64);index"`       // Unique across all apps
	WeChatConfig  string `gorm:"type:jsonb;default:'{}'"`      // JSON config: briefing time, etc.
	
	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	Memberships     []OrganizationMember `gorm:"foreignKey:UserID"`
	TeamMemberships []TeamMember         `gorm:"foreignKey:UserID"`
}
