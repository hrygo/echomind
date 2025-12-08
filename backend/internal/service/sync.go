package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/event"
	"github.com/hrygo/echomind/internal/repository"
	"github.com/hrygo/echomind/internal/tasks"
	"github.com/hrygo/echomind/pkg/event/bus"
	"github.com/hrygo/echomind/pkg/imap"
	echologger "github.com/hrygo/echomind/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var (
	ErrAccountNotConfigured = errors.New("email account not configured")
	syncTracer              = otel.Tracer("sync-service")
)

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
	accountRepo    repository.AccountRepository
	connector      IMAPConnector
	ingestor       *EmailIngestor
	bus            *bus.Bus         // Event Bus
	accountService *AccountService  // New dependency for account management
	config         *configs.Config  // Need full config to access security.EncryptionKey
	logger         CompatibleLogger // Add logger (兼容层)
}

// Ensure SyncService implements the EmailSyncer interface
var _ tasks.EmailSyncer = (*SyncService)(nil)

// NewSyncService creates a new SyncService.
func NewSyncService(accountRepo repository.AccountRepository, connector IMAPConnector, ingestor *EmailIngestor, bus *bus.Bus, accountService *AccountService, config *configs.Config, log echologger.Logger) *SyncService {
	return &SyncService{
		accountRepo:    accountRepo,
		connector:      connector,
		ingestor:       ingestor,
		bus:            bus,
		accountService: accountService,
		config:         config,
		logger:         echologger.AsZapSugaredLogger(log),
	}
}

// SyncEmails fetches emails for a specific user, saves them, and enqueues analysis tasks.
func (s *SyncService) SyncEmails(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, organizationID *uuid.UUID) error {
	ctx, span := syncTracer.Start(ctx, "SyncService.SyncEmails",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	if teamID != nil {
		span.SetAttributes(attribute.String("team.id", teamID.String()))
	}
	if organizationID != nil {
		span.SetAttributes(attribute.String("organization.id", organizationID.String()))
	}

	// 1. Get user's email account configuration (or team/org account)
	ctx, accountSpan := syncTracer.Start(ctx, "fetch_account_config")
	account, err := s.accountRepo.FindConfiguredAccount(ctx, userID, teamID, organizationID)
	accountSpan.End()

	if err != nil {
		span.RecordError(err)
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

	span.SetAttributes(
		attribute.String("account.email", account.Email),
		attribute.String("account.server", account.ServerAddress),
	)

	// 2. Connect to IMAP
	ctx, imapSpan := syncTracer.Start(ctx, "imap_connect")
	session, err := s.connector.Connect(ctx, account)
	imapSpan.End()

	if err != nil {
		imapSpan.RecordError(err)
		s.logger.Errorw("Failed to connect to IMAP", "error", err)
		return err
	}
	defer func() { _ = session.Logout() }()

	// 3. Ingest Emails
	// For now, we use a fixed time or logic. In future, use account.LastSyncAt
	lastSyncTime := time.Now().Add(-24 * time.Hour) // Default to 24h ago if not set
	if account.LastSyncAt != nil {
		lastSyncTime = *account.LastSyncAt
	}

	span.SetAttributes(
		attribute.String("sync.last_sync_time", lastSyncTime.Format(time.RFC3339)),
	)

	ctx, ingestSpan := syncTracer.Start(ctx, "ingest_emails")
	newEmails, err := s.ingestor.Ingest(ctx, session, account, lastSyncTime)
	ingestSpan.End()

	if err != nil {
		ingestSpan.RecordError(err)
		s.logger.Errorw("Failed to ingest emails", "error", err)
		return err
	}

	span.SetAttributes(
		attribute.Int("sync.new_emails_count", len(newEmails)),
	)

	// 4. Publish Events
	ctx, eventSpan := syncTracer.Start(ctx, "publish_events")
	for _, email := range newEmails {
		// Publish Email Synced Event
		event := event.EmailSyncedEvent{
			UserID: userID,
			Email:  email,
		}
		if err := s.bus.Publish(ctx, event); err != nil {
			s.logger.Errorw("Failed to publish email synced event",
				"email_id", email.ID,
				"user_id", userID,
				"error", err)
		}
	}
	eventSpan.End()

	return nil
}

// SyncEmailsForTask implements the EmailSyncer interface for use in background tasks
// It calls the main SyncEmails method with nil teamID and organizationID
func (s *SyncService) SyncEmailsForTask(ctx context.Context, userID uuid.UUID) error {
	return s.SyncEmails(ctx, userID, nil, nil)
}
