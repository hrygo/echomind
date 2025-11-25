package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/repository"
	echologger "github.com/hrygo/echomind/pkg/logger"
)

// EmailIngestor handles fetching and persisting emails.
type EmailIngestor struct {
	emailRepo repository.EmailRepository
	logger    CompatibleLogger
}

func NewEmailIngestor(emailRepo repository.EmailRepository, logger echologger.Logger) *EmailIngestor {
	return &EmailIngestor{
		emailRepo: emailRepo,
		logger:    echologger.AsZapSugaredLogger(logger),
	}
}

// Ingest fetches emails since the last sync time and saves them to the repository.
// It returns the list of newly saved emails.
func (s *EmailIngestor) Ingest(ctx context.Context, session IMAPSession, account *model.EmailAccount, lastSyncTime time.Time) ([]model.Email, error) {
	// 1. Fetch emails
	// For simplicity, we fetch a fixed limit for now, or we could use UID search based on lastSyncTime
	// The current fetcher implementation might need adjustment to support "since" properly if not already.
	// Assuming FetchEmails gets recent emails.
	emailDataList, err := session.FetchEmails("INBOX", 10) // Fetch top 10 for now
	if err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	var newEmails []model.Email
	userID := *account.UserID

	for _, data := range emailDataList {
		// Skip if email is older than lastSyncTime (if we want strict incremental sync)
		if data.Date.Before(lastSyncTime) {
			continue
		}

		// Check if email already exists
		exists, err := s.emailRepo.Exists(ctx, userID, data.MessageID)
		if err != nil {
			s.logger.Errorw("Failed to check email existence", "error", err)
			continue
		}
		if exists {
			continue
		}

		// Save new email
		email := model.Email{
			ID:        uuid.New(),
			UserID:    userID,
			AccountID: account.ID,
			Subject:   data.Subject,
			Sender:    data.Sender,
			Date:      data.Date,
			BodyText:  data.BodyText,
			BodyHTML:  data.BodyHTML,
			MessageID: data.MessageID,
			IsRead:    false, // Default to unread
			Folder:    "INBOX",
		}

		if err := s.emailRepo.Create(ctx, &email); err != nil {
			s.logger.Errorw("Failed to save email", "error", err)
			continue
		}

		newEmails = append(newEmails, email)
	}

	return newEmails, nil
}
