package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/imap"
	echologger "github.com/hrygo/echomind/pkg/logger"
	"github.com/hrygo/echomind/pkg/utils"
	"gorm.io/gorm"
)

var ErrAccountNotConfigured = errors.New("email account not configured")

type EmailFetcher interface {
	FetchEmails(c *clientimap.Client, mailbox string, limit int) ([]imap.EmailData, error)
}

type DefaultFetcher struct{}

func (d *DefaultFetcher) FetchEmails(c *clientimap.Client, mailbox string, limit int) ([]imap.EmailData, error) {
	return imap.FetchEmails(c, mailbox, limit)
}

// SyncService handles the business logic for synchronizing emails.

// IMAPClient defines the interface for IMAP client operations that SyncService needs.
type IMAPClient interface {
	DialAndLogin(addr, username, password string) (*clientimap.Client, error)
	Close(c *clientimap.Client)
}

// DefaultIMAPClient is the default implementation of IMAPClient using go-imap/client.
type DefaultIMAPClient struct{}

func (d *DefaultIMAPClient) DialAndLogin(addr, username, password string) (*clientimap.Client, error) {
	c, err := clientimap.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Login(username, password); err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

func (d *DefaultIMAPClient) Close(c *clientimap.Client) {
	c.Close()
}

// AsynqClientInterface defines the interface for asynq.Client operations that SyncService needs.
type AsynqClientInterface interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

// SyncService handles the business logic for synchronizing emails.
// CompatibleLogger 兼容的日志接口
type CompatibleLogger interface {
	Errorw(msg string, keysAndValues ...interface{})
	Errorf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

type SyncService struct {
	db             *gorm.DB
	imapClient     IMAPClient // Use the new IMAPClient interface
	fetcher        EmailFetcher
	asynqClient    AsynqClientInterface // Use the new AsynqClientInterface
	contactService *ContactService
	accountService *AccountService  // New dependency for account management
	config         *configs.Config  // Need full config to access security.EncryptionKey
	logger         CompatibleLogger // Add logger (兼容层)
}

// Ensure SyncService implements the EmailSyncer interface
var _ tasks.EmailSyncer = (*SyncService)(nil)

// NewSyncService creates a new SyncService.
func NewSyncService(db *gorm.DB, imapClient IMAPClient, fetcher EmailFetcher, asynqClient AsynqClientInterface, contactService *ContactService, accountService *AccountService, config *configs.Config, log echologger.Logger) *SyncService {
	return &SyncService{
		db:             db,
		imapClient:     imapClient,
		fetcher:        fetcher,
		asynqClient:    asynqClient,
		contactService: contactService,
		accountService: accountService,
		config:         config,
		logger:         echologger.AsZapSugaredLogger(log),
	}
}

// SyncEmails fetches emails for a specific user, saves them, and enqueues analysis tasks.
func (s *SyncService) SyncEmails(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, organizationID *uuid.UUID) error {
	// 1. Get user's email account configuration (or team/org account)
	var account model.EmailAccount
	query := s.db.WithContext(ctx).Model(&model.EmailAccount{})

	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	} else if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	} else {
		// Fallback to user-specific account if no team/org is provided
		query = query.Where("user_id = ?", userID)
	}

	err := query.First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "record not found") {
			return ErrAccountNotConfigured
		}
		s.logger.Errorw("Failed to fetch email account",
			"user_id", userID,
			"team_id", teamID,
			"org_id", organizationID,
			"error", err)
		return fmt.Errorf("failed to retrieve email account: %w", err)
	}

	// 2. Decrypt password
	keyBytes, err := hex.DecodeString(s.config.Security.EncryptionKey)
	if err != nil {
		return fmt.Errorf("invalid encryption key configuration: %w", err)
	}

	password, err := utils.Decrypt(account.EncryptedPassword, keyBytes)
	if err != nil {
		s.logger.Errorw("Error decrypting password for account %s: %v", account.ID, err)
		// Update account status to indicate decryption failure
		if statusErr := s.accountService.UpdateAccountStatus(ctx, account.ID, false, "Failed to decrypt password", nil); statusErr != nil {
			s.logger.Errorw("Failed to update account status after decryption failure: %v", statusErr)
		}
		return fmt.Errorf("failed to decrypt password: %w", err)
	}
	// 3. Establish IMAP connection (dynamic, per-sync)
	addr := fmt.Sprintf("%s:%d", account.ServerAddress, account.ServerPort)
	imapClient, err := s.imapClient.DialAndLogin(addr, account.Username, password)
	if err != nil {
		s.logger.Errorw("IMAP connection failed",
			"address", addr,
			"account_id", account.ID,
			"error", err)
		// Update account status to indicate connection failure
		if err := s.accountService.UpdateAccountStatus(ctx, account.ID, false, fmt.Sprintf("IMAP connection/login failed: %v", err), nil); err != nil {
			s.logger.Errorw("Failed to update account status after IMAP connection failure: %v", err)
		}
		return fmt.Errorf("failed to dial/login to IMAP server: %w", err)
	}
	defer s.imapClient.Close(imapClient)

	// Connection successful, update account status
	now := time.Now()
	if err := s.accountService.UpdateAccountStatus(ctx, account.ID, true, "", &now); err != nil {
		s.logger.Errorw("Failed to update account status after successful sync: %v", err)
	}
	// 4. Fetch emails using the dynamically created client
	emails, err := s.fetcher.FetchEmails(imapClient, "INBOX", 30)
	if err != nil {
		s.logger.Errorw("Failed to fetch emails",
			"account_id", account.ID,
			"error", err)
		// Optionally update account status with fetch error, keep connected true as login was successful
		if statusErr := s.accountService.UpdateAccountStatus(ctx, account.ID, true, fmt.Sprintf("Failed to fetch emails: %v", err), nil); statusErr != nil {
			s.logger.Errorw("Failed to update account status after email fetch failure: %v", statusErr)
		}
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
			email.ID = uuid.New() // Generate new UUID
			if err := s.db.WithContext(ctx).Create(&email).Error; err != nil {
				s.logger.Errorw("Failed to create email",
					"user_id", userID,
					"message_id", email.MessageID,
					"error", err)
				continue
			}
		} else {
			s.logger.Errorw("Database error",
				"user_id", userID,
				"error", result.Error)
			continue
		}

		// Enqueue Analysis Task
		if s.asynqClient != nil {
			task, err := tasks.NewEmailAnalyzeTask(email.ID, userID) // Pass UserID to task
			if err != nil {
				s.logger.Errorw("Failed to create analysis task",
					"email_id", email.ID,
					"user_id", userID,
					"error", err)
				continue
			}
			if _, err := s.asynqClient.Enqueue(task); err != nil {
				s.logger.Errorw("Failed to enqueue analysis task",
					"email_id", email.ID,
					"user_id", userID,
					"error", err)
			} else {
				s.logger.Debugw("Enqueued analysis task",
					"email_id", email.ID,
					"user_id", userID)
			}
		}

		// Update Contact Info (for this user)
		if s.contactService != nil && email.Sender != "" { // Ensure sender is not empty
			// Extract name from sender string if needed, or pass as is
			senderEmail, senderName := parseSender(email.Sender) // Assuming parseSender exists or create a simple one
			if senderEmail != "" {
				if err := s.contactService.UpdateContactFromEmail(ctx, userID, senderEmail, senderName, email.Date); err != nil {
					s.logger.Warnw("Failed to update contact",
						"user_id", userID,
						"email", senderEmail,
						"error", err)
				}
			}
		}
	}

	return nil
}

// SyncEmailsForTask implements the EmailSyncer interface for use in background tasks
// It calls the main SyncEmails method with nil teamID and organizationID
func (s *SyncService) SyncEmailsForTask(ctx context.Context, userID uuid.UUID) error {
	return s.SyncEmails(ctx, userID, nil, nil)
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
