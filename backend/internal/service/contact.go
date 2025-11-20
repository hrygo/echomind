package service

import (
	"time"

	"echomind.com/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{db: db}
}

// UpdateContactFromEmail extracts sender info and updates the contact record.
func (s *ContactService) UpdateContactFromEmail(email, name string, interactionTime time.Time) error {
	if email == "" {
		return nil
	}

	// Upsert contact:
	// If exists, increment count and update last interacted time.
	// If new, create.
	
	// Note: GORM's Upsert support varies by DB. For Postgres:
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":               name,
			"interaction_count":  gorm.Expr("interaction_count + 1"), // Removed "contact." prefix
			"last_interacted_at": interactionTime,
		}),
	}).Create(&model.Contact{
		Email:            email,
		Name:             name,
		InteractionCount: 1,
		LastInteractedAt: interactionTime,
	}).Error
}
