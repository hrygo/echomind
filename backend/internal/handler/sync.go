package handler

import (
	"net/http"

	"echomind.com/backend/internal/service"
	"echomind.com/backend/pkg/imap"
	clientimap "github.com/emersion/go-imap/client"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// EmailFetcher defines the interface for fetching email headers. 
type EmailFetcher interface {
	FetchEmails(c *clientimap.Client, mailbox string, limit int) ([]imap.EmailData, error)
}

type SyncHandler struct {
	db          *gorm.DB
	client      *clientimap.Client
	fetcher     EmailFetcher
	asynqClient *asynq.Client
}

// NewSyncHandler creates a new SyncHandler.
func NewSyncHandler(db *gorm.DB, imapClient *clientimap.Client, fetcher EmailFetcher, asynqClient *asynq.Client) *SyncHandler {
	return &SyncHandler{
		db:          db,
		client:      imapClient,
		fetcher:     fetcher,
		asynqClient: asynqClient,
	}
}

// SyncEmails handles the email synchronization request.
func (h *SyncHandler) SyncEmails(c *gin.Context) {
	// In a real application, this would trigger an async task.
	// For MVP, we'll run it directly (temporarily).

	err := service.SyncEmails(h.db, h.client, h.fetcher, h.asynqClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync initiated successfully"})
}
