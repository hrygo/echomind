package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type EmailEmbedding struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	EmailID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"email_id"`
	Content    string          `gorm:"type:text" json:"content"`                // Text chunk
	Vector     pgvector.Vector `gorm:"type:vector(1024)" json:"vector"`         // Fixed dimension: 1024
	Dimensions int             `gorm:"not null;default:1024" json:"dimensions"` // Always 1024
	CreatedAt  time.Time       `json:"created_at"`

	// Associations
	Email Email `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (EmailEmbedding) TableName() string {
	return "email_embeddings"
}
