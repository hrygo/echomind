package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

// EmailRepository defines the interface for email data access.
type EmailRepository interface {
	// Create saves a new email to the database.
	Create(ctx context.Context, email *model.Email) error
	// FindByMessageIDAndUserID finds an email by its Message-ID and UserID.
	FindByMessageIDAndUserID(ctx context.Context, messageID string, userID uuid.UUID) (*model.Email, error)
	// Save updates an existing email.
	Save(ctx context.Context, email *model.Email) error
	// Exists checks if an email exists by Message-ID and UserID.
	Exists(ctx context.Context, userID uuid.UUID, messageID string) (bool, error)
}

// GormEmailRepository is the GORM implementation of EmailRepository.
type GormEmailRepository struct {
	db *gorm.DB
}

// NewEmailRepository creates a new GormEmailRepository.
func NewEmailRepository(db *gorm.DB) *GormEmailRepository {
	return &GormEmailRepository{db: db}
}

// Create saves a new email to the database.
func (r *GormEmailRepository) Create(ctx context.Context, email *model.Email) error {
	return r.db.WithContext(ctx).Create(email).Error
}

// FindByMessageIDAndUserID finds an email by its Message-ID and UserID.
func (r *GormEmailRepository) FindByMessageIDAndUserID(ctx context.Context, messageID string, userID uuid.UUID) (*model.Email, error) {
	var email model.Email
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		First(&email).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// Save updates an existing email.
func (r *GormEmailRepository) Save(ctx context.Context, email *model.Email) error {
	return r.db.WithContext(ctx).Save(email).Error
}

// Exists checks if an email exists by Message-ID and UserID.
func (r *GormEmailRepository) Exists(ctx context.Context, userID uuid.UUID, messageID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Email{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
