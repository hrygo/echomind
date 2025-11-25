package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

type ActionService struct {
	db *gorm.DB
}

func NewActionService(db *gorm.DB) *ActionService {
	return &ActionService{db: db}
}

// ApproveEmail marks an email as approved/processed and archives it.
// In a real scenario, this might trigger an email reply or workflow.
func (s *ActionService) ApproveEmail(ctx context.Context, userID, emailID uuid.UUID) error {
	// 1. Verify ownership
	var email model.Email
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", emailID, userID).First(&email).Error; err != nil {
		return err
	}

	// 2. Perform Action: For now, we archive it (soft delete or move to folder?).
	// EchoMind "Approve" usually implies "Done".
	// Let's create a "Done" status or just archive it.
	// Since we don't have a status field other than folders (inferred), let's use a "Processed" tag or archive.
	// The plan says "archives email".
	// Let's soft delete it for now, or maybe we need an "Archived" flag?
	// The standard way in GORM for "Archive" without losing data is often just soft delete if that's the semantics,
	// or setting a Folder = "Archive" if we had that.
	// Current implementation of "trash" checks `deleted_at`.
	// Let's assume "Archive" means removing from Inbox but keeping it searchable.
	// The current `ListEmails` filters by `deleted_at IS NULL` by default.
	// So soft deleting removes it from Inbox.
	// Wait, "Trash" usually means soft delete in GORM defaults.
	// Maybe we need an `IsArchived` bool?
	// For now, let's use Soft Delete as "Archive/Done" to clear the inbox.
	// Or better: Add `IsArchived` field?
	// The plan said: "Approve/Dismiss".
	// Let's stick to Soft Delete for "Done" to get it out of the way.

	return s.db.WithContext(ctx).Delete(&email).Error
}

// SnoozeEmail hides the email until a specific time.
func (s *ActionService) SnoozeEmail(ctx context.Context, userID, emailID uuid.UUID, until time.Time) error {
	result := s.db.WithContext(ctx).Model(&model.Email{}).
		Where("id = ? AND user_id = ?", emailID, userID).
		Update("snoozed_until", until)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DismissEmail removes the email from the Smart Feed (High Urgency -> Low).
func (s *ActionService) DismissEmail(ctx context.Context, userID, emailID uuid.UUID) error {
	result := s.db.WithContext(ctx).Model(&model.Email{}).
		Where("id = ? AND user_id = ?", emailID, userID).
		Update("urgency", "Low")

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
