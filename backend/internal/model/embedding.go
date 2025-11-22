package model

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type EmailEmbedding struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	EmailID   uint             `gorm:"not null;index" json:"email_id"`
	Vector    pgvector.Vector  `gorm:"type:vector(1536)" json:"vector"` // 1536 for OpenAI text-embedding-3-small
	CreatedAt time.Time        `json:"created_at"`
	
	// Associations
	Email Email `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (EmailEmbedding) TableName() string {
	return "email_embeddings"
}
