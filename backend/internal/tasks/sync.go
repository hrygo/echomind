package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/pkg/logger"
)

const (
	TypeEmailSync = "email:sync"
)

type EmailSyncPayload struct {
	UserID uuid.UUID
}

// EmailSyncer defines the interface for syncing emails.
type EmailSyncer interface {
	SyncEmails(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, organizationID *uuid.UUID) error
}

// NewEmailSyncTask creates a task to sync emails for a user.
func NewEmailSyncTask(userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailSyncPayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailSync, payload), nil
}

// HandleEmailSyncTask processes the email sync task.
func HandleEmailSyncTask(ctx context.Context, t *asynq.Task, syncService EmailSyncer, log logger.Logger) error {
	var p EmailSyncPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.InfoContext(ctx, "Starting email sync",
		logger.String("user_id", p.UserID.String()))
	if err := syncService.SyncEmails(ctx, p.UserID, nil, nil); err != nil {
		log.ErrorContext(ctx, "Email sync failed",
			logger.String("user_id", p.UserID.String()),
			logger.Error(err))
		return fmt.Errorf("sync failed: %w", err)
	}
	log.InfoContext(ctx, "Email sync completed",
		logger.String("user_id", p.UserID.String()))
	return nil
}
