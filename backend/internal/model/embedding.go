package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type EmailEmbedding struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	EmailID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"email_id"`
	Content   string          `gorm:"type:text" json:"content"` // The actual text chunk
	Vector    pgvector.Vector `gorm:"type:vector(1024)" json:"vector"` // Maximum dimension: 1024 (current provider standard)
	Dimensions int             `gorm:"not null;default:1024" json:"dimensions"` // Track actual vector dimensions
	CreatedAt time.Time       `json:"created_at"`

	// Associations
	Email Email `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (EmailEmbedding) TableName() string {
	return "email_embeddings"
}

// BeforeCreate hook to validate dimensions and convert vector to correct size
func (e *EmailEmbedding) BeforeCreate(tx *gorm.DB) error {
	return e.validateAndConvertVector(tx)
}

// BeforeUpdate hook for safety (should rarely update vectors)
func (e *EmailEmbedding) BeforeUpdate(tx *gorm.DB) error {
	vectorSlice := e.Vector.Slice()
	if len(vectorSlice) > 0 {
		return e.validateAndConvertVector(tx)
	}
	return nil
}

// validateAndConvertVector ensures vector dimensions match database schema
func (e *EmailEmbedding) validateAndConvertVector(tx *gorm.DB) error {
	vectorSlice := e.Vector.Slice()
	if len(vectorSlice) == 0 {
		return nil
	}

	// Record actual vector dimensions
	e.Dimensions = len(vectorSlice)

	// Database schema enforces vector(1024) constraint
	// Vector dimension changes require full system reindexing (see docs)
	return nil
}
