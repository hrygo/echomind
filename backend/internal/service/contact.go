package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "email"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":               name,
			"interaction_count":  gorm.Expr("interaction_count + 1"),
			"last_interacted_at": interactionTime,
		}),
	}).Create(&model.Contact{
		ID:               uuid.New(), // Explicitly set ID
		UserID:           userID,
		Email:            email,
		Name:             name,
		InteractionCount: 1,
		LastInteractedAt: interactionTime,
	}).Error
}
