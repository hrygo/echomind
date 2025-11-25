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
	Vector    pgvector.Vector `gorm:"type:vector(1536)" json:"vector"` // Maximum dimension: 1536 (OpenAI standard)
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

// validateAndConvertVector validates vector dimensions and converts to max dimension if needed
func (e *EmailEmbedding) validateAndConvertVector(tx *gorm.DB) error {
	vectorSlice := e.Vector.Slice()
	if len(vectorSlice) == 0 {
		return nil
	}

	actualDimensions := len(vectorSlice)
	e.Dimensions = actualDimensions

	// If vector is longer than max dimension (1536), truncate it
	maxDimensions := 1536
	if actualDimensions > maxDimensions {
		truncatedSlice := vectorSlice[:maxDimensions]
		e.Vector = pgvector.NewVector(truncatedSlice)
		e.Dimensions = maxDimensions
	}

	// If vector is shorter than max dimensions, pad with zeros
	if actualDimensions < maxDimensions {
		paddedVector := make([]float32, maxDimensions)
		copy(paddedVector, vectorSlice)
		// Remaining elements remain zero
		e.Vector = pgvector.NewVector(paddedVector)
		// Keep original dimensions in the field
	}

	return nil
}
