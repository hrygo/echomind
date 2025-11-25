package event

import (
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
)

const EmailSyncedEventName = "email.synced"

// EmailSyncedEvent is published when an email is successfully synced.
type EmailSyncedEvent struct {
	UserID uuid.UUID
	Email  model.Email
}

func (e EmailSyncedEvent) Name() string {
	return EmailSyncedEventName
}
