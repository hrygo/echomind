package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{db: db}
}

// UpdateContactFromEmail extracts sender info and updates the contact record for a specific user.
func (s *ContactService) UpdateContactFromEmail(ctx context.Context, userID uuid.UUID, email, name string, interactionTime time.Time) error {
	if email == "" {
		return nil
	}

	// Upsert contact:
	// If exists, increment count and update last interacted time.
	// If new, create.

	// Note: GORM's Upsert support varies by DB. For Postgres:
	return s.db.WithContext(ctx).Create(&model.Contact{
		ID:               uuid.New(), // Explicitly set ID
		UserID:           &userID,
		Email:            email,
		Name:             name,
		InteractionCount: 1,
		LastInteractedAt: interactionTime,
	}).Error
}
