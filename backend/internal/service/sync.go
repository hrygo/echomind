package service

import (
	"log"

	"echomind.com/backend/internal/model"
	"echomind.com/backend/internal/tasks"
	"echomind.com/backend/pkg/imap"
	"github.com/emersion/go-imap/client"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type EmailFetcher interface {
	FetchEmails(c *client.Client, mailbox string, limit int) ([]imap.EmailData, error)
}

type DefaultFetcher struct{}

func (d *DefaultFetcher) FetchEmails(c *client.Client, mailbox string, limit int) ([]imap.EmailData, error) {
	return imap.FetchEmails(c, mailbox, limit)
}

// SyncEmails fetches emails using the provided fetcher, saves them, and enqueues analysis tasks.
func SyncEmails(db *gorm.DB, c *client.Client, fetcher EmailFetcher, asynqClient *asynq.Client) error {
	// Fetch latest 30 emails
	emails, err := fetcher.FetchEmails(c, "INBOX", 30)
	if err != nil {
		return err
	}

	for _, h := range emails {
		email := model.Email{
			MessageID: h.MessageID,
			Subject:   h.Subject,
			Sender:    h.Sender,
			Date:      h.Date,
			Snippet:   "",
			BodyText:  h.BodyText,
			BodyHTML:  h.BodyHTML,
		}

		// Upsert: On Conflict Do Nothing. 
		// We need to get the ID if it exists or is created to enqueue task.
		// GORM's OnConflict doesn't easily return ID if existing.
		// Strategy: Try Create. If conflict, First search by MessageID.
		
		var savedEmail model.Email
		
		// Check if exists
		result := db.Where("message_id = ?", email.MessageID).First(&savedEmail)
		if result.Error == nil {
			// Exists. Skip if already summarized? 
			// For MVP, we might want to re-analyze or skip. Let's skip if Summary exists.
			if savedEmail.Summary != "" {
				continue
			}
			// Update ID for task
			email.ID = savedEmail.ID
		} else if result.Error == gorm.ErrRecordNotFound {
			// Create new
			if err := db.Create(&email).Error; err != nil {
				log.Printf("Failed to create email: %v", err)
				continue
			}
		} else {
			log.Printf("DB error: %v", result.Error)
			continue
		}
		
		// Enqueue Analysis Task
		if asynqClient != nil {
			task, err := tasks.NewEmailAnalyzeTask(email.ID)
			if err != nil {
				log.Printf("Failed to create task for email %d: %v", email.ID, err)
				continue
			}
			if _, err := asynqClient.Enqueue(task); err != nil {
				log.Printf("Failed to enqueue task for email %d: %v", email.ID, err)
			} else {
                log.Printf("Enqueued analysis task for email %d", email.ID)
            }
		}
	}

	return nil
}
