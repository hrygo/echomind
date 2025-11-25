package listener

import (
	"context"
	"fmt"
	"strings"

	"github.com/hrygo/echomind/internal/event"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/event/bus"
	echologger "github.com/hrygo/echomind/pkg/logger"
)

// CompatibleLogger defines the interface for structured logging
type CompatibleLogger interface {
	Errorw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

// AnalysisListener handles email analysis tasks.
type AnalysisListener struct {
	asynqClient service.AsynqClientInterface
	logger      CompatibleLogger
}

func NewAnalysisListener(client service.AsynqClientInterface, logger echologger.Logger) *AnalysisListener {
	return &AnalysisListener{
		asynqClient: client,
		logger:      echologger.AsZapSugaredLogger(logger),
	}
}

func (l *AnalysisListener) Handle(ctx context.Context, e bus.Event) error {
	evt, ok := e.(event.EmailSyncedEvent)
	if !ok {
		return fmt.Errorf("invalid event type: %T", e)
	}

	if l.asynqClient == nil {
		return nil
	}

	task, err := tasks.NewEmailAnalyzeTask(evt.Email.ID, evt.UserID)
	if err != nil {
		l.logger.Errorw("Failed to create analysis task",
			"email_id", evt.Email.ID,
			"user_id", evt.UserID,
			"error", err)
		return err
	}

	if _, err := l.asynqClient.Enqueue(task); err != nil {
		l.logger.Errorw("Failed to enqueue analysis task",
			"email_id", evt.Email.ID,
			"user_id", evt.UserID,
			"error", err)
		return err
	}

	l.logger.Debugw("Enqueued analysis task",
		"email_id", evt.Email.ID,
		"user_id", evt.UserID)
	return nil
}

// ContactListener handles contact updates from emails.
type ContactListener struct {
	contactService *service.ContactService
	logger         CompatibleLogger
}

func NewContactListener(contactService *service.ContactService, logger echologger.Logger) *ContactListener {
	return &ContactListener{
		contactService: contactService,
		logger:         echologger.AsZapSugaredLogger(logger),
	}
}

func (l *ContactListener) Handle(ctx context.Context, e bus.Event) error {
	evt, ok := e.(event.EmailSyncedEvent)
	if !ok {
		return fmt.Errorf("invalid event type: %T", e)
	}

	if l.contactService == nil || evt.Email.Sender == "" {
		return nil
	}

	senderEmail, senderName := parseSender(evt.Email.Sender)
	if senderEmail != "" {
		if err := l.contactService.UpdateContactFromEmail(ctx, evt.UserID, senderEmail, senderName, evt.Email.Date); err != nil {
			l.logger.Warnw("Failed to update contact",
				"user_id", evt.UserID,
				"email", senderEmail,
				"error", err)
			return err
		}
	}
	return nil
}

// Helper to parse sender string (duplicated from SyncService for now, could be moved to utils)
func parseSender(sender string) (email, name string) {
	if idx := strings.LastIndex(sender, "<"); idx != -1 {
		if endIdx := strings.LastIndex(sender, ">"); endIdx != -1 && endIdx > idx {
			email = sender[idx+1 : endIdx]
			name = strings.TrimSpace(sender[:idx])
			return
		}
	}
	email = sender
	return
}
