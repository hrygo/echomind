package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type EmailEmbedding struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	EmailID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"email_id"`
	Vector    pgvector.Vector `gorm:"type:vector(1536)" json:"vector"` // 1536 for OpenAI text-embedding-3-small
	CreatedAt time.Time       `json:"created_at"`

	// Associations
	Email Email `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (EmailEmbedding) TableName() string {
	return "email_embeddings"
}
