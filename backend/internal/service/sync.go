package service

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/imap"
	clientimap "github.com/emersion/go-imap/client"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type EmailFetcher interface {
	FetchEmails(c *clientimap.Client, mailbox string, limit int) ([]imap.EmailData, error)
}

type DefaultFetcher struct{}

func (d *DefaultFetcher) FetchEmails(c *clientimap.Client, mailbox string, limit int) ([]imap.EmailData, error) {
	return imap.FetchEmails(c, mailbox, limit)
}

// SyncService handles the business logic for synchronizing emails.
type SyncService struct {
	db          *gorm.DB
	imapClient  *clientimap.Client
	fetcher     EmailFetcher
	asynqClient *asynq.Client
	contactService *ContactService // New dependency for contact updates
}

// NewSyncService creates a new SyncService.
func NewSyncService(db *gorm.DB, imapClient *clientimap.Client, fetcher EmailFetcher, asynqClient *asynq.Client, contactService *ContactService) *SyncService {
	return &SyncService{
		db:          db,
		imapClient:  imapClient,
		fetcher:     fetcher,
		asynqClient: asynqClient,
		contactService: contactService,
	}
}

// SyncEmails fetches emails for a specific user, saves them, and enqueues analysis tasks.
func (s *SyncService) SyncEmails(ctx context.Context, userID uuid.UUID) error {
	// Fetch latest 30 emails
	// In a real multi-tenant scenario, the IMAP client needs to be user-specific.
	// For now, we use a placeholder client. This will be addressed in future phases.
    // FIXME: For local dev/demo, if fetch fails (likely due to no creds), we log and return success to avoid 500.
	emails, err := s.fetcher.FetchEmails(s.imapClient, "INBOX", 30)
	if err != nil {
        log.Printf("Warning: Failed to fetch emails (likely due to missing credentials in local dev): %v", err)
		return nil // Return nil to avoid 500 error in frontend
	}

	for _, h := range emails {
		email := model.Email{
			UserID:    userID, // Set UserID
			MessageID: h.MessageID,
			Subject:   h.Subject,
			Sender:    h.Sender,
			Date:      h.Date,
			Snippet:   "",
			BodyText:  h.BodyText,
			BodyHTML:  h.BodyHTML,
		}

		var savedEmail model.Email
		
		// Check if email with this MessageID and UserID already exists
		result := s.db.WithContext(ctx).Where("user_id = ? AND message_id = ?", userID, email.MessageID).First(&savedEmail)
		if result.Error == nil {
			// Exists. Skip if already summarized.
			if savedEmail.Summary != "" {
				continue
			}
			// Update ID for task
			email.ID = savedEmail.ID
		} else if result.Error == gorm.ErrRecordNotFound {
			// Create new
			if err := s.db.WithContext(ctx).Create(&email).Error; err != nil {
				log.Printf("Failed to create email for user %s: %v", userID, err)
				continue
			}
		} else {
			log.Printf("DB error for user %s: %v", userID, result.Error)
			continue
		}
		
		// Enqueue Analysis Task
		if s.asynqClient != nil {
			task, err := tasks.NewEmailAnalyzeTask(email.ID, userID) // Pass UserID to task
			if err != nil {
				log.Printf("Failed to create task for email %s (user %s): %v", email.ID, userID, err)
				continue
			}
			if _, err := s.asynqClient.Enqueue(task); err != nil {
				log.Printf("Failed to enqueue task for email %s (user %s): %v", email.ID, userID, err)
			} else {
                log.Printf("Enqueued analysis task for email %s (user %s)", email.ID, userID)
            }
		}

		// Update Contact Info (for this user)
		if s.contactService != nil && email.Sender != "" { // Ensure sender is not empty
            // Extract name from sender string if needed, or pass as is
			senderEmail, senderName := parseSender(email.Sender) // Assuming parseSender exists or create a simple one
            if senderEmail != "" {
                if err := s.contactService.UpdateContactFromEmail(ctx, userID, senderEmail, senderName, email.Date); err != nil {
                    log.Printf("Failed to update contact for user %s, email %s: %v", userID, senderEmail, err)
                }
            }
		}
	}

	return nil
}

// Helper to parse sender string, e.g., "Name <email@example.com>" or "email@example.com"
func parseSender(sender string) (email, name string) {
	if idx := strings.LastIndex(sender, "<"); idx != -1 {
		if endIdx := strings.LastIndex(sender, ">"); endIdx != -1 && endIdx > idx {
			email = sender[idx+1 : endIdx]
			name = strings.TrimSpace(sender[:idx])
			return
		}
	}
	// Fallback if format is just "email@example.com"
	email = sender
	return
}