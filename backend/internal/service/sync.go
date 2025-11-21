package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/imap"
	clientimap "github.com/emersion/go-imap/client"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/pkg/utils"
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
	fetcher     EmailFetcher
	asynqClient *asynq.Client
	contactService *ContactService
	accountService *AccountService // New dependency for account management
	config *configs.Config // Need full config to access security.EncryptionKey
}

// NewSyncService creates a new SyncService.
func NewSyncService(db *gorm.DB, fetcher EmailFetcher, asynqClient *asynq.Client, contactService *ContactService, accountService *AccountService, config *configs.Config) *SyncService {
	return &SyncService{
		db:          db,
		fetcher:     fetcher,
		asynqClient: asynqClient,
		contactService: contactService,
		accountService: accountService,
		config: config,
	}
}

// SyncEmails fetches emails for a specific user, saves them, and enqueues analysis tasks.
func (s *SyncService) SyncEmails(ctx context.Context, userID uuid.UUID) error {
	// 1. Get user's email account configuration
	account, err := s.accountService.GetAccountByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error fetching email account for user %s: %v", userID, err)
		return fmt.Errorf("failed to retrieve email account: %w", err)
	}

	// 2. Decrypt password
	decryptionKey := s.config.Security.EncryptionKey // Access config via accountService
	password, err := utils.Decrypt(account.EncryptedPassword, []byte(decryptionKey))
	if err != nil {
		log.Printf("Error decrypting password for account %s: %v", account.ID, err)
		// Update account status to indicate decryption failure
		s.accountService.UpdateAccountStatus(ctx, account.ID, false, "Failed to decrypt password", nil)
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// 3. Establish IMAP connection (dynamic, per-sync)
	addr := fmt.Sprintf("%s:%d", account.ServerAddress, account.ServerPort)
	imapClient, err := clientimap.DialTLS(addr, nil)
	if err != nil {
		log.Printf("Error dialing IMAP server %s for account %s: %v", addr, account.ID, err)
		// Update account status to indicate connection failure
		s.accountService.UpdateAccountStatus(ctx, account.ID, false, fmt.Sprintf("IMAP connection failed: %v", err), nil)
		return fmt.Errorf("failed to dial IMAP server: %w", err)
	}
	defer imapClient.Close()

	if err := imapClient.Login(account.Username, password); err != nil {
		log.Printf("Error logging into IMAP for account %s: %v", account.ID, err)
		// Update account status to indicate login failure
		s.accountService.UpdateAccountStatus(ctx, account.ID, false, fmt.Sprintf("IMAP login failed: %v", err), nil)
		return fmt.Errorf("failed to login to IMAP server: %w", err)
	}

	// Connection successful, update account status
	now := time.Now()
	s.accountService.UpdateAccountStatus(ctx, account.ID, true, "", &now)

	// 4. Fetch emails using the dynamically created client
	emails, err := s.fetcher.FetchEmails(imapClient, "INBOX", 30)
	if err != nil {
		log.Printf("Error fetching emails for account %s: %v", account.ID, err)
		// Optionally update account status with fetch error, keep connected true as login was successful
		s.accountService.UpdateAccountStatus(ctx, account.ID, true, fmt.Sprintf("Failed to fetch emails: %v", err), nil)
		return fmt.Errorf("failed to fetch emails: %w", err)
	}

	for _, h := range emails {
		email := model.Email{
			UserID:    userID,
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