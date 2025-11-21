package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hrygo/echomind/internal/model"
)

// EmailService handles email-related business logic and data access.
type EmailService struct {
	db *gorm.DB
}

// NewEmailService creates a new EmailService.
func NewEmailService(db *gorm.DB) *EmailService {
	return &EmailService{
		db: db,
	}
}

// ListEmails retrieves a list of emails for a given user.
func (s *EmailService) ListEmails(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Email, error) {
	var emails []model.Email
	query := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}

// GetEmail retrieves a single email by its ID for a given user.
func (s *EmailService) GetEmail(ctx context.Context, userID, emailID uuid.UUID) (*model.Email, error) {
	var email model.Email
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&email, emailID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Email not found or not owned by user
		}
		return nil, err
	}
	return &email, nil
}

// CreateEmail creates a new email record.
func (s *EmailService) CreateEmail(ctx context.Context, email *model.Email) error {
	return s.db.WithContext(ctx).Create(email).Error
}

// UpdateEmail updates an existing email record.
func (s *EmailService) UpdateEmail(ctx context.Context, email *model.Email) error {
	return s.db.WithContext(ctx).Save(email).Error
}
