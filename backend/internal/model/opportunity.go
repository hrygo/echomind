package model

import (
	"time"

	"gorm.io/gorm"
)

// OpportunityType represents the type of opportunity
type OpportunityType string

const (
	OpportunityTypeBuying      OpportunityType = "buying"
	OpportunityTypePartnership OpportunityType = "partnership"
	OpportunityTypeRenewal     OpportunityType = "renewal"
	OpportunityTypeStrategic   OpportunityType = "strategic"
)

// OpportunityStatus represents the status of an opportunity
type OpportunityStatus string

const (
	OpportunityStatusNew       OpportunityStatus = "new"
	OpportunityStatusActive    OpportunityStatus = "active"
	OpportunityStatusWon       OpportunityStatus = "won"
	OpportunityStatusLost      OpportunityStatus = "lost"
	OpportunityStatusOnHold    OpportunityStatus = "on_hold"
)

// Opportunity represents a business opportunity
type Opportunity struct {
	ID          string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title       string           `gorm:"not null;size:500" json:"title"`
	Description string           `gorm:"type:text" json:"description"`
	Company     string           `gorm:"not null;size:200" json:"company"`
	Value       string           `gorm:"size:100" json:"value"`
	Type        OpportunityType  `gorm:"type:opportunity_type;default:'buying'" json:"type"`
	Status      OpportunityStatus `gorm:"type:opportunity_status;default:'new'" json:"status"`
	Confidence  int              `gorm:"check:confidence >= 0 AND confidence <= 100" json:"confidence"`
	UserID      string           `gorm:"type:uuid;not null;index" json:"user_id"`
	TeamID      string           `gorm:"type:uuid;index" json:"team_id"`
	OrgID       string           `gorm:"type:uuid;index" json:"org_id"`
	SourceEmailID *string         `gorm:"type:uuid" json:"source_email_id"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`

	// Associations
	Contacts   []Contact   `gorm:"many2many:opportunity_contacts;" json:"contacts,omitempty"`
	Activities []Activity `gorm:"foreignKey:OpportunityID" json:"activities,omitempty"`
}

// TableName returns the table name for the Opportunity model
func (Opportunity) TableName() string {
	return "opportunities"
}

// OpportunityContact represents the relationship between opportunities and contacts
type OpportunityContact struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OpportunityID string    `gorm:"type:uuid;not null"`
	ContactID     string    `gorm:"type:uuid;not null"`
	Role          string    `gorm:"size:100"` // e.g., "decision maker", "influencer", "technical contact"
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the OpportunityContact model
func (OpportunityContact) TableName() string {
	return "opportunity_contacts"
}

// Activity represents activities related to an opportunity
type Activity struct {
	ID            string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OpportunityID string `gorm:"type:uuid;not null;index"`
	UserID        string `gorm:"type:uuid;not null"`
	Type          string `gorm:"size:50"` // e.g., "call", "email", "meeting", "note"
	Title         string `gorm:"not null;size:200"`
	Description   string `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the Activity model
func (Activity) TableName() string {
	return "activities"
}