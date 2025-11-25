package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

// AccountRepository defines the interface for email account data access.
type AccountRepository interface {
	// FindConfiguredAccount finds the email account configuration for a user, team, or organization.
	// Priority: Organization > Team > User (based on provided IDs)
	FindConfiguredAccount(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, organizationID *uuid.UUID) (*model.EmailAccount, error)
}

// GormAccountRepository is the GORM implementation of AccountRepository.
type GormAccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new GormAccountRepository.
func NewAccountRepository(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{db: db}
}

// FindConfiguredAccount finds the email account configuration.
func (r *GormAccountRepository) FindConfiguredAccount(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, organizationID *uuid.UUID) (*model.EmailAccount, error) {
	var account model.EmailAccount
	query := r.db.WithContext(ctx).Model(&model.EmailAccount{})

	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	} else if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	} else {
		// Fallback to user-specific account if no team/org is provided
		query = query.Where("user_id = ?", userID)
	}

	err := query.First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
