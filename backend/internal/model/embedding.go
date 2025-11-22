package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type EmailEmbedding struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	EmailID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"email_id"`
	Content   string          `gorm:"type:text" json:"content"`       // The actual text chunk
	Vector    pgvector.Vector `gorm:"type:vector(768)" json:"vector"` // 768 for Gemini text-embedding-004
	CreatedAt time.Time       `json:"created_at"`

	// Associations
	Email Email `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (EmailEmbedding) TableName() string {
	return "email_embeddings"
}
