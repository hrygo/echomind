package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeEmailSync = "email:sync"
)

type EmailSyncPayload struct {
	UserID uuid.UUID
}

// EmailSyncer defines the interface for syncing emails.
type EmailSyncer interface {
	SyncEmails(ctx context.Context, userID uuid.UUID) error
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
func HandleEmailSyncTask(ctx context.Context, t *asynq.Task, syncService EmailSyncer) error {
	var p EmailSyncPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Starting email sync for user: %s", p.UserID)
	if err := syncService.SyncEmails(ctx, p.UserID); err != nil {
		log.Printf("Email sync failed for user %s: %v", p.UserID, err)
		return fmt.Errorf("sync failed: %w", err)
	}
	log.Printf("Email sync completed for user: %s", p.UserID)
	return nil
}
